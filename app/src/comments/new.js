import style from './index.module.scss';
import React, {Component} from 'react';
import API from '../api/index.js';
import Button from '../components/button.js';

class CommentNew extends Component {
  constructor(props) {
    super(props);

    this.api = new API();
    this.state = {
      topic_id: props.topicId,
      body: '',
      submitting: false
    }

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(e) {
    const {name, value} = e.target;
    this.setState({
      [name]: value
    });
  }

  handleSubmit(e) {
    e.preventDefault();
    if (this.state.submitting) {
      return
    }
    this.setState({submitting: true});
    this.api.comment.create(this.state).then((resp) => {
      if (resp.error) {
        this.setState({submitting: false});
        return
      }
      this.props.handleSubmit(resp.data);
      this.setState({body: '', submitting: false});
    });
  }

  render() {
    const i18n = window.i18n;
    let state = this.state;
    if (!this.api.user.loggedIn()) {
      return (
        <div className={style.custom}>
          {i18n.t('comment.custom')}
        </div>
      )
    }
    return (
      <div className={style.form}>
        <form onSubmit={this.handleSubmit}>
          <input type='hidden' name='topic_id' defaultValue={state.topic_id} />
          <div>
            <textarea type='text' name='body' minLength='3' required placeholder={i18n.t('comment.form.body')} value={state.body} onChange={this.handleChange} />
          </div>
          <div className='action'>
            <Button type='submit' classes='submit' disabled={state.submitting} text={i18n.t('general.submit')}/>
          </div>
        </form>
      </div>
    )
  }
}

export default CommentNew;
