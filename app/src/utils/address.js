export const shortAddress = (address) => {
  if (address.length !== 42) {
    return address;
  }
  return `${address.substring(0, 4)}...${address.slice(-4, address.length)}`;
};
