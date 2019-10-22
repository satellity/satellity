class ColorUtils {
  constructor() {
    this.colors = ['#EF564F', '#4B93D1', '#9354CA', '#414141'];
  }

  colour(str) {
    let i = Math.floor(Math.random() * 100);
    if (str !== '') {
      i = str.charCodeAt(0);
    }
    return this.colors[i % 3];
  }
}

export default ColorUtils;
