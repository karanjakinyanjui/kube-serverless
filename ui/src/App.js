import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import './App.css';
import FunctionList from './components/FunctionList';
import FunctionDetail from './components/FunctionDetail';
import DeployFunction from './components/DeployFunction';
import Dashboard from './components/Dashboard';

function App() {
  return (
    <Router>
      <div className="App">
        <nav className="navbar">
          <div className="nav-brand">
            <h1>Kube-Serverless</h1>
          </div>
          <div className="nav-links">
            <Link to="/">Dashboard</Link>
            <Link to="/functions">Functions</Link>
            <Link to="/deploy">Deploy</Link>
          </div>
        </nav>

        <div className="container">
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/functions" element={<FunctionList />} />
            <Route path="/functions/:name" element={<FunctionDetail />} />
            <Route path="/deploy" element={<DeployFunction />} />
          </Routes>
        </div>
      </div>
    </Router>
  );
}

export default App;
