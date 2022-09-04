export const shortAddress = (address) => {
  return `${address.substring(0, 4)}...${address.slice(-4, address.length)}`;
};
