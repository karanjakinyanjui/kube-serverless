const express = require('express');
const bodyParser = require('body-parser');
const fs = require('fs');
const path = require('path');

const app = express();
const PORT = process.env.PORT || 8080;

app.use(bodyParser.json());

let handler;
let coldStart = true;

// Load function code
const loadFunction = () => {
  try {
    const codePath = '/function/code';
    const handlerName = process.env.FUNCTION_HANDLER || 'index.handler';

    if (fs.existsSync(codePath)) {
      const code = fs.readFileSync(codePath, 'utf8');

      // Create a temporary module
      const modulePath = path.join('/tmp', 'function.js');
      fs.writeFileSync(modulePath, code);

      // Load the module
      const fn = require(modulePath);

      // Parse handler name (e.g., "index.handler" -> "handler")
      const parts = handlerName.split('.');
      const exportName = parts[parts.length - 1];

      handler = fn[exportName] || fn.handler || fn;

      console.log('Function loaded successfully');
      coldStart = false;
    } else {
      console.warn('No function code found, using echo handler');
      handler = async (event) => {
        return {
          statusCode: 200,
          body: JSON.stringify({ echo: event })
        };
      };
    }
  } catch (error) {
    console.error('Error loading function:', error);
    handler = async (event) => {
      return {
        statusCode: 500,
        body: JSON.stringify({ error: error.message })
      };
    };
  }
};

// Initialize function
loadFunction();

// Health check
app.get('/health', (req, res) => {
  res.json({ status: 'healthy' });
});

// Ready check
app.get('/ready', (req, res) => {
  res.json({ status: 'ready', coldStart });
});

// Function invocation
app.post('/', async (req, res) => {
  const startTime = Date.now();
  const wasColdStart = coldStart;

  if (coldStart) {
    coldStart = false;
  }

  try {
    const event = {
      body: req.body,
      headers: req.headers,
      method: req.method,
      path: req.path,
      query: req.query
    };

    const result = await handler(event);

    const duration = Date.now() - startTime;

    res.set('X-Function-Duration', duration.toString());
    res.set('X-Cold-Start', wasColdStart.toString());

    if (result && typeof result === 'object' && result.statusCode) {
      res.status(result.statusCode).send(result.body);
    } else {
      res.json(result);
    }
  } catch (error) {
    console.error('Error executing function:', error);
    res.status(500).json({ error: error.message });
  }
});

// Metrics endpoint
app.get('/metrics', (req, res) => {
  res.set('Content-Type', 'text/plain');
  res.send(`# HELP function_cold_start Cold start indicator
# TYPE function_cold_start gauge
function_cold_start ${coldStart ? 1 : 0}
`);
});

app.listen(PORT, () => {
  console.log(`Node.js runtime server listening on port ${PORT}`);
  console.log(`Function: ${process.env.FUNCTION_NAME}`);
  console.log(`Handler: ${process.env.FUNCTION_HANDLER}`);
});
