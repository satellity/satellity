const colors = ['#EF564F', '#4B93D1', '#9354CA', '#877D6A', '#F9BA48'];

export const colorful = (str, alpha) => {
  let i = Math.floor(Math.random() * 100);
  if (str !== '') {
    i = str.charCodeAt(0);
  }
  alpha = Math.round(alpha * 255);
  const hex = (alpha + 0x10000).toString(16).substr(-2).toUpperCase();
  return colors[i % 5] + hex;
};
