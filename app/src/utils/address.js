export const shortAddress = (address) => {
  return `${address.substring(0, 8)}...${address.slice(-8, address.length)}`;
};
