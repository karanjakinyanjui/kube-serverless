import json
from datetime import datetime

def handler(event):
    """
    Process data from event and return summary
    """
    print(f"Processing data at {datetime.now().isoformat()}")

    # Extract data from event
    data = event.get('body', {})

    # Perform data processing (example)
    processed = {
        'processed_at': datetime.now().isoformat(),
        'record_count': len(data) if isinstance(data, list) else 1,
        'summary': 'Data processed successfully'
    }

    return {
        'statusCode': 200,
        'body': processed
    }
