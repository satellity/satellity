class Me {
  constructor(api) {
    this.api = api;
  }

  signOut() {
    window.localStorage.clear();
  }
}

export default Me;
