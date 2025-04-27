import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import './App.css'
import RegistrationPage from './pages/RegistrationPage'
import LoginPage from './pages/LoginPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<RegistrationPage />} />
        <Route path="/login" element={<LoginPage />} />
      </Routes>
    </Router>
  );
}

export default App
