import React from 'react';
import '@testing-library/jest-dom/extend-expect';
import { render } from '@testing-library/react';
import CookbookApp from '../CookbookApp';

test('renders learn react link', () => {
  const { getByText } = render(<CookbookApp />);
  const linkElement = getByText(/Cookbook/i);
  expect(linkElement).toBeInTheDocument();
});
