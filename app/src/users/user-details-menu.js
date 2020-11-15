import React, { Component } from "react";
import { Redirect } from "react-router-dom";
import { Dropdown, DropdownButton } from "react-bootstrap";
import API from "../api/index.js";
import Avatar from "./avatar.js";
import "./user-details-menu.scss";

class UserDetailsMenu extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    const user = this.api.user.local();
    this.state = {
      user: user,
      nickname: user.nickname,
      me: this.api.user.loggedIn(),
    };
  }

  handleSignOut = (e) => {
    console.log("signout");
    e.preventDefault();

    this.api.me.signOut();
    this.setState({ me: false });
  };

  render() {
    if (!this.state.me) {
      return <Redirect to={{ pathname: "/" }} />;
    }

    return (
      <DropdownButton
        menuAlign={{ sm: "right" }}
        title={<Avatar user={this.state.user} class="small" />}
        variant="link"
        id="user-details-dropdown-menu"
        className="user-details-dropdown-menu-component"
      >
        <Dropdown.Item eventKey="1" href="/user/edit">
          Signed in as <strong>{this.state.nickname}</strong>
        </Dropdown.Item>
        <Dropdown.Divider />
        <Dropdown.Item eventKey="2" onClick={this.handleSignOut}>
          {window.i18n.t("general.sign.out")}
        </Dropdown.Item>
      </DropdownButton>
    );
  }
}

export default UserDetailsMenu;
