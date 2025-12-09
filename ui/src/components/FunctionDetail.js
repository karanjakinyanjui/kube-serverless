import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

function FunctionDetail() {
  const { name } = useParams();
  const navigate = useNavigate();
  const [functionData, setFunctionData] = useState(null);
  const [metrics, setMetrics] = useState(null);
  const [loading, setLoading] = useState(true);
  const [invokePayload, setInvokePayload] = useState('{}');
  const [invokeResult, setInvokeResult] = useState(null);

  useEffect(() => {
    fetchFunction();
    fetchMetrics();
  }, [name]);

  const fetchFunction = async () => {
    try {
      const response = await axios.get(`${API_URL}/api/v1/functions/${name}`);
      setFunctionData(response.data);
    } catch (error) {
      console.error('Error fetching function:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchMetrics = async () => {
    try {
      const response = await axios.get(`${API_URL}/api/v1/functions/${name}/metrics`);
      setMetrics(response.data);
    } catch (error) {
      console.error('Error fetching metrics:', error);
    }
  };

  const invokeFunction = async () => {
    try {
      const payload = JSON.parse(invokePayload);
      const response = await axios.post(`${API_URL}/api/v1/functions/${name}/invoke`, payload);
      setInvokeResult(JSON.stringify(response.data, null, 2));
    } catch (error) {
      setInvokeResult(`Error: ${error.message}`);
    }
  };

  const deleteFunction = async () => {
    if (!window.confirm(`Are you sure you want to delete function "${name}"?`)) {
      return;
    }

    try {
      await axios.delete(`${API_URL}/api/v1/functions/${name}`);
      navigate('/functions');
    } catch (error) {
      console.error('Error deleting function:', error);
      alert('Failed to delete function');
    }
  };

  if (loading) {
    return <div className="loading">Loading...</div>;
  }

  if (!functionData) {
    return <div className="card">Function not found</div>;
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
        <h2>{functionData.name}</h2>
        <button onClick={deleteFunction} className="btn btn-danger">
          Delete Function
        </button>
      </div>

      <div className="card">
        <h3>Details</h3>
        <table>
          <tbody>
            <tr>
              <th>Runtime</th>
              <td>{functionData.runtime}</td>
            </tr>
            <tr>
              <th>Handler</th>
              <td>{functionData.handler}</td>
            </tr>
            <tr>
              <th>Min Replicas</th>
              <td>{functionData.minReplicas || 0}</td>
            </tr>
            <tr>
              <th>Max Replicas</th>
              <td>{functionData.maxReplicas || 10}</td>
            </tr>
            <tr>
              <th>Current Replicas</th>
              <td>{functionData.status?.replicas || 0}</td>
            </tr>
            <tr>
              <th>Status</th>
              <td>
                <span className={`status-badge status-${functionData.status?.state === 'running' ? 'running' : 'error'}`}>
                  {functionData.status?.state || 'unknown'}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      {metrics && (
        <div className="card">
          <h3>Metrics</h3>
          <div className="metrics-grid">
            <div className="metric-card">
              <div className="metric-label">Invocations</div>
              <div className="metric-value">{metrics.invocations}</div>
            </div>
            <div className="metric-card">
              <div className="metric-label">Cold Starts</div>
              <div className="metric-value">{metrics.coldStarts}</div>
            </div>
            <div className="metric-card">
              <div className="metric-label">Avg Duration</div>
              <div className="metric-value">{metrics.avgDuration.toFixed(3)}s</div>
            </div>
            <div className="metric-card">
              <div className="metric-label">Error Rate</div>
              <div className="metric-value">{(metrics.errorRate * 100).toFixed(2)}%</div>
            </div>
            <div className="metric-card">
              <div className="metric-label">Cost Estimate</div>
              <div className="metric-value">${metrics.costEstimate.toFixed(4)}</div>
            </div>
          </div>
        </div>
      )}

      <div className="card">
        <h3>Test Invocation</h3>
        <div className="form-group">
          <label>Payload (JSON)</label>
          <textarea
            value={invokePayload}
            onChange={(e) => setInvokePayload(e.target.value)}
            placeholder='{"key": "value"}'
          />
        </div>
        <button onClick={invokeFunction} className="btn btn-primary">
          Invoke Function
        </button>
        {invokeResult && (
          <div style={{ marginTop: '1rem' }}>
            <h4>Result:</h4>
            <pre style={{ background: '#f8f9fa', padding: '1rem', borderRadius: '4px', overflow: 'auto' }}>
              {invokeResult}
            </pre>
          </div>
        )}
      </div>
    </div>
  );
}

export default FunctionDetail;
