import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

function Dashboard() {
  const [functions, setFunctions] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchFunctions();
  }, []);

  const fetchFunctions = async () => {
    try {
      const response = await axios.get(`${API_URL}/api/v1/functions`);
      setFunctions(response.data || []);
    } catch (error) {
      console.error('Error fetching functions:', error);
    } finally {
      setLoading(false);
    }
  };

  const totalFunctions = functions.length;
  const runningFunctions = functions.filter(f => f.status?.state === 'running').length;
  const totalReplicas = functions.reduce((sum, f) => sum + (f.status?.replicas || 0), 0);

  if (loading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div>
      <h2>Dashboard</h2>

      <div className="metrics-grid">
        <div className="metric-card">
          <div className="metric-label">Total Functions</div>
          <div className="metric-value">{totalFunctions}</div>
        </div>
        <div className="metric-card">
          <div className="metric-label">Running Functions</div>
          <div className="metric-value">{runningFunctions}</div>
        </div>
        <div className="metric-card">
          <div className="metric-label">Total Replicas</div>
          <div className="metric-value">{totalReplicas}</div>
        </div>
        <div className="metric-card">
          <div className="metric-label">Platform Status</div>
          <div className="metric-value" style={{ fontSize: '1.5rem', color: '#2ecc71' }}>
            Healthy
          </div>
        </div>
      </div>

      <div className="card">
        <h3>Platform Overview</h3>
        <p>
          Kube-Serverless is a Function-as-a-Service (FaaS) platform running on Kubernetes.
          Deploy serverless functions with auto-scaling, event-driven triggers, and comprehensive monitoring.
        </p>
      </div>

      <div className="card">
        <h3>Recent Functions</h3>
        {functions.length === 0 ? (
          <p>No functions deployed yet. Deploy your first function to get started!</p>
        ) : (
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Runtime</th>
                <th>Replicas</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {functions.slice(0, 5).map(fn => (
                <tr key={fn.name}>
                  <td>{fn.name}</td>
                  <td>{fn.runtime}</td>
                  <td>{fn.status?.replicas || 0}</td>
                  <td>
                    <span className={`status-badge status-${fn.status?.state === 'running' ? 'running' : 'error'}`}>
                      {fn.status?.state || 'unknown'}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}

export default Dashboard;
