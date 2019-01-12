import './loading.scss';
import React from 'react';

const LoadingView = (props) => (
  <div className={`lds-ring ${props.style}`}>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
  </div>
);

export default LoadingView;
