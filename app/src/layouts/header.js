import React, {useState, useEffect} from 'react';
import {Link} from 'react-router-dom';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {ethers} from 'ethers';
import {generateKeyPair} from '@stablelib/ed25519';
import {encode} from '@stablelib/hex';
import Web3Modal from 'web3modal';
import API from 'api/index.js';
import Config from 'components/config.js';
import {shortAddress} from 'utils';

import logo from 'assets/images/logo.svg';
import style from './main.module.scss';

const providerOptions = {
  /* See Provider Options Section */
};

const Header = () => {
  const api = new API();
  const user = api.user;
  const meData = api.me;

  const [address, setAddress] = useState('');
  const [me, setMe] = useState(meData.value());
  const [web3Modal, setWeb3Modal] = useState(null);

  const handleLoginClick = async (e) => {
    if (address !== '') {
      return;
    }
    web3Modal.clearCachedProvider();
    const instance = await web3Modal.connect();
    const provider = new ethers.providers.Web3Provider(instance);

    // Subscribe to accounts change
    instance.on('accountsChanged', (accounts) => {
      console.log('accountsChanged', accounts);
    });

    // Subscribe to chainId change
    instance.on('chainChanged', (chainId) => {
      console.log('chainChanged', chainId);
    });

    // Subscribe to instance connection
    instance.on('connect', (info) => {
      console.log('connect', info);
    });

    // Subscribe to instance disconnection
    instance.on('disconnect', (error) => {
      // console.log('disconnect', error);
    });
    const userAddress = await provider.getSigner().getAddress();
    const key = generateKeyPair();
    const sessionPublic = encode(key.publicKey, true);
    const sessionPrivate = encode(key.secretKey, true);
    const msg = ethers.utils.id(`Satellite::${userAddress}:${sessionPublic}`);
    const sig = await provider.getSigner().signMessage(msg);
    user.create(userAddress, sessionPublic, sessionPrivate, sig.slice(2)).then((resp) => {
      if (resp.error) {
        return;
      }
      setMe(resp.data);
    });
    setAddress(userAddress);
  };

  useEffect(() => {
    const newWeb3Modal = new Web3Modal({
      cacheProvider: true, // optional
      providerOptions, // required
    });
    setWeb3Modal(newWeb3Modal);
  }, []);

  useEffect(() => {
    if (web3Modal && web3Modal.cachedProvider) {
      handleLoginClick();
    }
  }, [web3Modal]);

  let profile = <span className={style.navi} onClick={handleLoginClick}>Login</span>;
  if (!!me) {
    profile = (
      <div className={style.navis}>
        <Link to='/topics/new' className={`${style.navi}`}> <FontAwesomeIcon icon={['fa', 'plus']} /> </Link>
        <Link to='/user/edit' className={`${style.navi} ${style.user}`}> {shortAddress(me.nickname)} </Link>
      </div>
    );
  }

  return (
    <header className={style.header}>
      <Link className={style.site} to='/'>
        <img className={style.logo} src={logo} alt={Config.Name} />
        <span className={style.name}>{Config.Name}</span>
      </Link>

      <div className={style.menus}>
        <Link className={`${style.menu} ${window.location.pathname === '/' ? style.current : ''}` } to='/'>
          Home
        </Link>
      </div>
      {profile}
    </header>
  );
};

export default Header;
