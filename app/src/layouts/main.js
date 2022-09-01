import React from 'react';
import {Outlet} from 'react-router-dom';
import Header from './header.js';

const MainLayout = () => (
  <>
    <Header />
    <div className='wrapper'>
      <Outlet />
    </div>
  </>
);

export default MainLayout;
