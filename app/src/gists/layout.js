import React from 'react';
import Chart from './chart.js';
import PropTypes from 'prop-types';
import Widget from 'components/widget.js';

const View = ({children}) => {
  return (
    <div className='container'>
      <main className='column main'>
        {children}
      </main>
      <aside className='column aside'>
        <Widget>
          <Chart />
        </Widget>
      </aside>
    </div>
  );
};

View.propTypes = {
  children: PropTypes.any,
};

export default View;
