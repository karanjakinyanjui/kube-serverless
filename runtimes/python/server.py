import os
import sys
import time
import json
import importlib.util
from flask import Flask, request, jsonify
from prometheus_client import Counter, Histogram, generate_latest, CONTENT_TYPE_LATEST

app = Flask(__name__)

cold_start = True
handler = None

# Metrics
invocations = Counter('function_invocations_total', 'Total function invocations')
duration = Histogram('function_duration_seconds', 'Function execution duration')
cold_starts = Counter('function_cold_starts_total', 'Total cold starts')

def load_function():
    """Load the function from the mounted code"""
    global handler, cold_start

    try:
        code_path = '/function/code'
        handler_name = os.getenv('FUNCTION_HANDLER', 'handler.handler')

        if os.path.exists(code_path):
            with open(code_path, 'r') as f:
                code = f.read()

            # Create a module from the code
            spec = importlib.util.spec_from_loader('function_module', loader=None)
            module = importlib.util.module_from_spec(spec)
            exec(code, module.__dict__)

            # Parse handler name (e.g., "handler.handler" -> "handler")
            parts = handler_name.split('.')
            export_name = parts[-1]

            handler = getattr(module, export_name, None)
            if handler is None:
                handler = getattr(module, 'handler', None)

            print(f'Function loaded successfully: {handler_name}')
            cold_start = False
        else:
            print('No function code found, using echo handler')
            handler = lambda event: {'statusCode': 200, 'body': json.dumps({'echo': event})}

    except Exception as e:
        print(f'Error loading function: {e}')
        handler = lambda event: {'statusCode': 500, 'body': json.dumps({'error': str(e)})}

# Initialize function
load_function()

@app.route('/health', methods=['GET'])
def health():
    return jsonify({'status': 'healthy'})

@app.route('/ready', methods=['GET'])
def ready():
    return jsonify({'status': 'ready', 'coldStart': cold_start})

@app.route('/', methods=['POST'])
def invoke():
    global cold_start

    start_time = time.time()
    was_cold_start = cold_start

    if cold_start:
        cold_starts.inc()
        cold_start = False

    invocations.inc()

    try:
        event = {
            'body': request.get_json(silent=True) or {},
            'headers': dict(request.headers),
            'method': request.method,
            'path': request.path,
            'query': dict(request.args)
        }

        result = handler(event)

        exec_duration = time.time() - start_time
        duration.observe(exec_duration)

        response = jsonify(result)
        response.headers['X-Function-Duration'] = str(exec_duration)
        response.headers['X-Cold-Start'] = str(was_cold_start)

        if isinstance(result, dict) and 'statusCode' in result:
            response.status_code = result['statusCode']

        return response

    except Exception as e:
        print(f'Error executing function: {e}')
        return jsonify({'error': str(e)}), 500

@app.route('/metrics', methods=['GET'])
def metrics():
    return generate_latest(), 200, {'Content-Type': CONTENT_TYPE_LATEST}

if __name__ == '__main__':
    port = int(os.getenv('PORT', 8080))
    print(f'Python runtime server listening on port {port}')
    print(f'Function: {os.getenv("FUNCTION_NAME")}')
    print(f'Handler: {os.getenv("FUNCTION_HANDLER")}')
    app.run(host='0.0.0.0', port=port)
