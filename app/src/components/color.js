class ColorUtils {

  constructor() {
    this.colors = ['#EF564F', '#4B93D1', '#9354CA', '#414141'];
  }

  colour(str, alpha) {
    let i = Math.floor(Math.random() * 100);
    if (str !== '') {
      i = str.charCodeAt(0);
    }
    alpha = Math.round(alpha * 255);
    let hex = (alpha + 0x10000).toString(16).substr(-2).toUpperCase();
    console.log(this.colors[i % 3] + hex)
    return this.colors[i % 3] + hex;
  }
}

export default ColorUtils;
