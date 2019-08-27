import style from './new.scss';
import React, {Component} from 'react';
import API from '../api/index.js';
import LoadingView from '../loading/loading.js';

class New extends Component {
  constructor(props) {
    super(props);
    this.state = {
      group_id: props.groupId,
      body: '',
      submitting: false
    }

    this.api = new API();
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDisMount() {}

  handleChange(e) {
    const target = e.target;
    const name = target.name;
    this.setState({
      [name]: target.value
    });
  }

  handleSubmit(e) {
    e.preventDefault();
    let state = this.state;
    this.setState({submitting: true}, () => {
      this.api.message.create(state.group_id, {body: state.body}).then((data) => {
        this.setState({submitting: false, body: ''});
      });
    });
  }

  render() {
    let state = this.state;
    return (
      <div className={style.form}>
        <form onSubmit={this.handleSubmit}>
          <input type='hidden' name='group_id' defaultValue={state.group_id} />
          <div>
            <textarea
              type='text'
              name='body'
              minLength='3'
              required
              value={state.body}
              onChange={this.handleChange} />
          </div>
          <div className='action'>
            <button className='btn submit' disabled={state.submitting}>
              { state.submitting && <LoadingView style='sm-ring blank'/> }
              &nbsp;{i18n.t('general.submit')}
            </button>
          </div>
        </form>
      </div>
    )
  }
}

export default New;
