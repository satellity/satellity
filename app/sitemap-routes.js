import React from 'react';
import { Route } from 'react-router';

export default (
  <Route>
    <Route path='/topics/:id' />
    <Route exact path='/products' />
    <Route path='/products/:id' />
    <Route path='/products/q/:id' />
  </Route>
);
