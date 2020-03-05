import React from 'react';

type State = {
  input: string;
  error: string;
  isPending: boolean;
};

type Action =
  // prettier-ignore
  | { type: 'UPDATE' } & Partial<Pick<State, 'input' | 'error'>>
  | { type: 'SEND' }
  | { type: 'DONE' };

const stateReducer = (prevState: State, action: Action) => {
  switch (action.type) {
    case 'UPDATE': {
      const { type, ...rest } = action;

      return {
        ...prevState,
        ...rest,
      };
    }

    case 'SEND': {
      return {
        ...prevState,
        error: '',
        isPending: true,
      };
    }

    case 'DONE': {
      return {
        ...prevState,
        isPending: false,
      };
    }

    default: {
      return prevState;
    }
  }
};

export const useSingleInput = (initialState: State) => {
  return React.useReducer(stateReducer, initialState);
};
