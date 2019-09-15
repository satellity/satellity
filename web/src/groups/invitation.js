import React, {Component} from 'react';
import API from '../api/index.js';

class Invitation extends Component {
  constructor(props) {
    super(props);

    this.state = {
      group_id: props.groupId,
      email: '',
      submitting: false
    }
    this.api = new API();
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(e) {
    const target = e.target;
    const name = target.name;
    this.setState({
      [name]: target.value
    });
  }

  handleSubmit(e) {
    e.preventDefault();
    if (this.state.submitting) {
      return
    }
    this.setState({submitting: true}, () => {
      const data = {email: this.state.email};
      this.api.group.invite(this.state.group_id, data).then(() => {
        this.setState({submitting: false, email:''});
      });
    });
  }

  render() {
    return (
      <div>
        <form onSubmit={this.handleSubmit}>
          <div>
            <input type='text' name='email' required value={this.state.email} autoComplete='off' placeholder='email' onChange={this.handleChange} />
          </div>
          <div>
            <button type="submit" className='btn invite' disabled={this.state.submitting}>
                &nbsp;{i18n.t('general.submit')}
            </button>
          </div>
        </form>
      </div>
    )
  }
}

export default Invitation;
