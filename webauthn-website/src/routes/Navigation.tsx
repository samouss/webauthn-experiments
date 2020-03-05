import React from 'react';
import { Link, useLocation } from 'react-router-dom';

type Props = {
  isAuthenticated: boolean;
  onLogout: () => void;
};

const Navigation: React.FC<Props> = ({ isAuthenticated, onLogout }) => {
  const location = useLocation();

  return (
    <nav className="breadcrumb has-bullet-separator" aria-label="breadcrumbs">
      <ul>
        {!isAuthenticated ? (
          <>
            <li className={location.pathname === '/' ? 'is-active' : ''}>
              <Link to="/">Register</Link>
            </li>
            <li className={location.pathname === '/login' ? 'is-active' : ''}>
              <Link to="/login">Login</Link>
            </li>
          </>
        ) : (
          <>
            <li className={location.pathname === '/me' ? 'is-active' : ''}>
              <Link to="/me">Home</Link>
            </li>
            <li>
              <a
                href="#"
                onClick={event => {
                  event.preventDefault();
                  onLogout();
                }}
              >
                Logout
              </a>
            </li>
          </>
        )}
      </ul>
    </nav>
  );
};

export default Navigation;
