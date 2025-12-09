import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

function FunctionList() {
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

  const deleteFunction = async (name) => {
    if (!window.confirm(`Are you sure you want to delete function "${name}"?`)) {
      return;
    }

    try {
      await axios.delete(`${API_URL}/api/v1/functions/${name}`);
      fetchFunctions();
    } catch (error) {
      console.error('Error deleting function:', error);
      alert('Failed to delete function');
    }
  };

  if (loading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
        <h2>Functions</h2>
        <Link to="/deploy" className="btn btn-primary">Deploy New Function</Link>
      </div>

      {functions.length === 0 ? (
        <div className="card">
          <p>No functions deployed yet.</p>
        </div>
      ) : (
        <div className="card">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Runtime</th>
                <th>Handler</th>
                <th>Replicas</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {functions.map(fn => (
                <tr key={fn.name}>
                  <td>
                    <Link to={`/functions/${fn.name}`} style={{ color: '#3498db', textDecoration: 'none' }}>
                      {fn.name}
                    </Link>
                  </td>
                  <td>{fn.runtime}</td>
                  <td>{fn.handler}</td>
                  <td>{fn.status?.replicas || 0}</td>
                  <td>
                    <span className={`status-badge status-${fn.status?.state === 'running' ? 'running' : 'error'}`}>
                      {fn.status?.state || 'unknown'}
                    </span>
                  </td>
                  <td>
                    <button
                      onClick={() => deleteFunction(fn.name)}
                      className="btn btn-danger"
                      style={{ fontSize: '0.85rem', padding: '0.25rem 0.75rem' }}
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

export default FunctionList;
