import {ethers} from 'ethers';
import Web3Modal from 'web3modal';

const providerOptions = {
  /* See Provider Options Section */
};

const web3Modal = new Web3Modal({
  cacheProvider: true, // optional
  providerOptions, // required
});

export const useProvider = async () => {
  const instance = await web3Modal.connect();
  return new ethers.providers.Web3Provider(instance);
};
