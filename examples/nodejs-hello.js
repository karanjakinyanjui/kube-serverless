module.exports.handler = async (event) => {
  console.log('Received event:', JSON.stringify(event, null, 2));

  return {
    statusCode: 200,
    body: JSON.stringify({
      message: 'Hello from Kube-Serverless!',
      timestamp: new Date().toISOString(),
      headers: event.headers
    })
  };
};
