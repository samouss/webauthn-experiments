import React from 'react';
import { BrowserRouter, Link, Switch, Route, Redirect } from 'react-router-dom';
import Register from './routes/Register';
import Login from './routes/Login';
import Home from './routes/Home';
import Navigation from './routes/Navigation';
import { createClient } from './http';
import './App.css';

const isWebauthnAvailable = 'PublicKeyCredential' in window;

type State = {
  isAuthenticated: boolean;
  token: string;
};

type Action =
  // prettier-ignore
  | { type: 'LOGIN', token: string }
  | { type: 'LOGOUT' };

const stateReducer = (prevState: State, action: Action) => {
  switch (action.type) {
    case 'LOGIN': {
      return {
        ...prevState,
        isAuthenticated: true,
        token: action.token,
      };
    }

    case 'LOGOUT': {
      return {
        ...prevState,
        isAuthenticated: false,
        token: '',
      };
    }

    default: {
      return prevState;
    }
  }
};

const App = () => {
  const [state, dispatch] = React.useReducer(stateReducer, {
    isAuthenticated: false,
    token: '',
  });

  const client = React.useMemo(
    () =>
      createClient({
        endpoint: '/api',
        token: state.token,
      }),
    [state.token]
  );

  return (
    <BrowserRouter>
      <section className="hero is-medium is-dark is-bold">
        <div className="hero-body">
          <div className="container">
            <h1 className="title">Webauthn</h1>
            <h2 className="subtitle">
              Passwordless authentication with Webauthn API.
            </h2>
          </div>
        </div>
      </section>

      <section className="section">
        <div className="container">
          {!isWebauthnAvailable && (
            <div className="notification is-danger">
              Your browser doesn't support{' '}
              <a href="https://developer.mozilla.org/en-US/docs/Web/API/Web_Authentication_API">
                Webauthn
              </a>{' '}
              API.
            </div>
          )}

          <div className="notification">
            This is an experiment using{' '}
            <a href="https://developer.mozilla.org/en-US/docs/Web/API/Web_Authentication_API">
              Webauthn
            </a>{' '}
            API. The data sent to the server is only stored in memory. We don't
            persist any information to disk or database. Each time the server
            restarts data is lost on purpose. The goal is to demonstrate how
            passwordless authentication is possible. You can take a look at the
            code for both the server/client on{' '}
            <a href="http://github.com/samouss/webauthn-experiments">GitHub</a>.
          </div>

          <Navigation
            isAuthenticated={state.isAuthenticated}
            onLogout={() => {
              dispatch({ type: 'LOGOUT' });
            }}
          />

          <Switch>
            <Route path="/me">
              {state.isAuthenticated ? (
                <Home client={client} />
              ) : (
                <Redirect to="/login" />
              )}
            </Route>

            <Route path="/login">
              {state.isAuthenticated ? (
                <Redirect to="/me" />
              ) : (
                <Login
                  enabled={isWebauthnAvailable}
                  client={client}
                  onLogin={token => {
                    dispatch({
                      type: 'LOGIN',
                      token,
                    });
                  }}
                />
              )}
            </Route>

            <Route path="/" exact>
              {state.isAuthenticated ? (
                <Redirect to="/me" />
              ) : (
                <Register
                  enabled={isWebauthnAvailable}
                  client={client}
                  onRegister={token => {
                    dispatch({
                      type: 'LOGIN',
                      token,
                    });
                  }}
                />
              )}
            </Route>

            <Route path="*">
              <div className="notification is-warning">
                This is not the web page you are looking for.{' '}
                <Link to="/">Back to the home page.</Link>
              </div>
            </Route>
          </Switch>
        </div>
      </section>
    </BrowserRouter>
  );
};

export default App;
