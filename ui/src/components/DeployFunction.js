import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

function DeployFunction() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    name: '',
    runtime: 'nodejs18',
    handler: 'index.handler',
    code: '',
    minReplicas: 0,
    maxReplicas: 10
  });

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      await axios.post(`${API_URL}/api/v1/functions`, formData);
      alert('Function deployed successfully!');
      navigate('/functions');
    } catch (error) {
      console.error('Error deploying function:', error);
      alert('Failed to deploy function: ' + (error.response?.data || error.message));
    }
  };

  const loadExample = (runtime) => {
    const examples = {
      nodejs18: `module.exports.handler = async (event) => {
  return {
    statusCode: 200,
    body: JSON.stringify({
      message: 'Hello from Node.js!',
      event: event
    })
  };
};`,
      python39: `def handler(event):
    return {
        'statusCode': 200,
        'body': {
            'message': 'Hello from Python!',
            'event': event
        }
    }`,
      go119: `package main

func Handler(event map[string]interface{}) (interface{}, error) {
    return map[string]interface{}{
        "statusCode": 200,
        "body": map[string]interface{}{
            "message": "Hello from Go!",
            "event":   event,
        },
    }, nil
}`
    };

    setFormData({
      ...formData,
      code: examples[runtime] || ''
    });
  };

  return (
    <div>
      <h2>Deploy Function</h2>

      <div className="card">
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Function Name *</label>
            <input
              type="text"
              name="name"
              value={formData.name}
              onChange={handleChange}
              required
              placeholder="my-function"
            />
          </div>

          <div className="form-group">
            <label>Runtime *</label>
            <select
              name="runtime"
              value={formData.runtime}
              onChange={handleChange}
              required
            >
              <option value="nodejs18">Node.js 18</option>
              <option value="python39">Python 3.9</option>
              <option value="go119">Go 1.19</option>
            </select>
            <button
              type="button"
              onClick={() => loadExample(formData.runtime)}
              className="btn btn-success"
              style={{ marginTop: '0.5rem', fontSize: '0.85rem', padding: '0.25rem 0.75rem' }}
            >
              Load Example
            </button>
          </div>

          <div className="form-group">
            <label>Handler *</label>
            <input
              type="text"
              name="handler"
              value={formData.handler}
              onChange={handleChange}
              required
              placeholder="index.handler"
            />
          </div>

          <div className="form-group">
            <label>Function Code *</label>
            <textarea
              name="code"
              value={formData.code}
              onChange={handleChange}
              required
              placeholder="Enter your function code here..."
            />
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
            <div className="form-group">
              <label>Min Replicas</label>
              <input
                type="number"
                name="minReplicas"
                value={formData.minReplicas}
                onChange={handleChange}
                min="0"
              />
            </div>

            <div className="form-group">
              <label>Max Replicas</label>
              <input
                type="number"
                name="maxReplicas"
                value={formData.maxReplicas}
                onChange={handleChange}
                min="1"
              />
            </div>
          </div>

          <button type="submit" className="btn btn-primary">
            Deploy Function
          </button>
        </form>
      </div>
    </div>
  );
}

export default DeployFunction;
