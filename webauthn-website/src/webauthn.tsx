import { Client } from './http';

type RegisterOptions = {
  client: Client;
  user: {
    id: string;
    email: string;
  };
};

type StartRegisterResponse = {
  publicKey: {
    challenge: string;
    pubKeyCredParams: PublicKeyCredentialParameters[];
    rp: PublicKeyCredentialRpEntity;
    user: PublicKeyCredentialUserEntity & {
      id: string;
    };
  };
};

type CompleteRegisterResponse = {
  token: string;
};

export const register = (options: RegisterOptions) =>
  options.client
    .post<StartRegisterResponse>('/register/start', options.user)
    .then(res => {
      const options = res.publicKey;

      return navigator.credentials.create({
        publicKey: {
          ...options,
          challenge: strToBuf(options.challenge),
          user: {
            ...options.user,
            id: strToBuf(options.user.id),
          },
        },
      });
    })
    .then(credential => {
      if (!credential) {
        throw new Error('Fail to create a new credential.');
      }

      if (!isPublicKeyCredential(credential)) {
        throw new Error('Fail to create the public key credential.');
      }

      if (!isAuthenticatorAttestationResponse(credential.response)) {
        throw new Error('Fail to create the public key credential.');
      }

      return options.client.post<CompleteRegisterResponse>(
        '/register/complete',
        {
          id: credential.id,
          rawId: bufToStr(credential.rawId),
          type: credential.type,
          response: {
            attestationObject: bufToStr(credential.response.attestationObject),
            clientDataJSON: bufToStr(credential.response.clientDataJSON),
          },
        }
      );
    });

type LoginOptions = {
  client: Client;
  user: {
    email: string;
  };
};

type StartLoginResponse = {
  publicKey: {
    challenge: string;
    allowCredentials: (PublicKeyCredentialDescriptor & {
      id: string;
    })[];
  };
};

type CompleteLoginResponse = {
  token: string;
};

export const login = (options: LoginOptions) =>
  options.client
    .post<StartLoginResponse>('/login/start', options.user)
    .then(res => {
      const options = res.publicKey;

      return navigator.credentials.get({
        publicKey: {
          ...options,
          challenge: strToBuf(options.challenge),
          allowCredentials: options.allowCredentials.map(credential => ({
            ...credential,
            id: strToBuf(credential.id),
          })),
        },
      });
    })
    .then(credential => {
      if (!credential) {
        throw new Error('Fail to get the credential.');
      }

      if (!isPublicKeyCredential(credential)) {
        throw new Error('Fail to get the public key credential.');
      }

      if (!isAuthenticatorAssertionResponse(credential.response)) {
        throw new Error('Fail to get the public key credential.');
      }

      return options.client.post<CompleteLoginResponse>('/login/complete', {
        id: credential.id,
        rawId: bufToStr(credential.rawId),
        type: credential.type,
        response: {
          authenticatorData: bufToStr(credential.response.authenticatorData),
          signature: bufToStr(credential.response.signature),
          clientDataJSON: bufToStr(credential.response.clientDataJSON),
          ...(credential.response.userHandle && {
            userHandle: bufToStr(credential.response.userHandle),
          }),
        },
      });
    });

export const isAbortableError = (x: Error): boolean =>
  x instanceof DOMException &&
  (x.name === 'NotAllowedError' || x.name === 'AbortError');

// PublicKeyCredential is identified by its type which is: public-key
// https://developer.mozilla.org/en-US/docs/Web/API/PublicKeyCredential
const isPublicKeyCredential = (x: Credential): x is PublicKeyCredential => {
  return x.type === 'public-key';
};

const isAuthenticatorAttestationResponse = (
  x: AuthenticatorResponse
): x is AuthenticatorAttestationResponse => {
  return x instanceof AuthenticatorAttestationResponse;
};

const isAuthenticatorAssertionResponse = (
  x: AuthenticatorResponse
): x is AuthenticatorAssertionResponse => {
  return x instanceof AuthenticatorResponse;
};

const strToBuf = (x: string) => Uint8Array.from(atob(x), c => c.charCodeAt(0));

const bufToStr = (buffer: ArrayBuffer) => {
  return (
    btoa(
      new Uint8Array(buffer).reduce(
        (acc, value) => (acc += String.fromCharCode(value)),
        ''
      )
    )
      // Go is not able to decode the base64 value with those characters.
      .replace(/\+/g, '-')
      .replace(/\//g, '_')
      .replace(/=/g, '')
  );
};
