import React from 'react';
import { Client } from '../http';
import { login, isAbortableError } from '../webauthn';
import { useSingleInput } from '../input';

type Props = {
  client: Client;
  enabled: boolean;
  onLogin: (token: string) => void;
};

const Login: React.FC<Props> = ({ client, enabled, onLogin }) => {
  const [state, dispatch] = useSingleInput({
    input: '',
    error: '',
    isPending: false,
  });

  return (
    <form
      action=""
      onSubmit={(event: React.FormEvent) => {
        event.preventDefault();
        const user = {
          email: state.input,
        };

        dispatch({ type: 'SEND' });
        login({ client, user })
          .then(response => {
            dispatch({ type: 'DONE' });
            onLogin(response.token);
          })
          .catch(response => {
            dispatch({ type: 'DONE' });

            // An error wrapper could help to differentiate login/register errors
            // vs navigator/fetch errors. We could display the former and report
            // the latter.
            if (isAbortableError(response)) {
              // This error might happen because the user cancel on purpose. The
              // navigator already display a message (on it's explicitly done by
              // the user). We don't have to remind him the action he took.
              return;
            }

            if (response instanceof Error) {
              // An error wrapper could help to differentiate login/register errors
              // vs navigator/fetch errors. We could display the former and report
              // the latter.
              return dispatch({
                type: 'UPDATE',
                error: 'An unexpected error has occurred.',
              });
            }

            return dispatch({
              type: 'UPDATE',
              error: response.message,
            });
          });
      }}
    >
      <div className="field">
        <div className="control">
          <input
            className={`input${state.error ? ' is-danger' : ''}`}
            type="email"
            placeholder="Email"
            value={state.input}
            onChange={event =>
              dispatch({
                type: 'UPDATE',
                input: event.currentTarget.value,
              })
            }
            required
          />
        </div>
        {state.error && <p className="help is-danger">{state.error}</p>}
      </div>

      <div className="field is-grouped is-grouped-right">
        <div className="control">
          <button
            className={`button is-link${state.isPending ? ' is-loading' : ''}`}
            disabled={!enabled}
          >
            Login
          </button>
        </div>
      </div>
    </form>
  );
};

export default Login;
