import style from './new.scss';
require('codemirror/lib/codemirror.css');
require('codemirror/theme/xq-light.css');
require('codemirror/mode/markdown/markdown.js');
import React, { Component } from 'react';
import { Redirect } from 'react-router-dom';
import {Controlled as CodeMirror} from 'react-codemirror2'
import API from '../api/index.js';
const validate = require('uuid-validate');

class New extends Component {
  constructor(props) {
    super(props);
    this.api = new API();

    let id = this.props.match.params.id;
    if (!id) {
      id = ''
    }
    this.state = {
      group_id: id,
      name: '',
      description: '',
      submitting: false,
      loading: false
    }
    this.handleChange = this.handleChange.bind(this);
    this.handleDescriptionChange = this.handleDescriptionChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDidMount() {
    if (validate(this.state.group_id)) {
      this.api.group.show(this.state.group_id).then((data) => {
        data.loading = false;
        this.setState(data);
      });
    }
  }

  handleChange(e) {
    const target = e.target;
    const name = target.name;
    this.setState({
      [name]: target.value
    });
  }

  handleDescriptionChange(editor, data, value) {
    this.setState({description: value});
  }

  handleSubmit(e) {
    e.preventDefault();
    if (this.state.submitting) {
      return
    }
    this.setState({submitting: true}, () => {
      const history = this.props.history;
      const data = {name: this.state.name, description: this.state.description};
      let request;
      if (validate(this.state.group_id)) {
        request = this.api.group.update(this.state.group_id, data);
      } else {
        request = this.api.group.create(data);
      };
      request.then((data) => {
        this.setState({submitting: false});
        // TODO should use other uri
        history.push('/');
      });
    });
  }

  render() {
    let state = this.state;

    if (!this.api.user.loggedIn()) {
      return (
        <Redirect to={{ pathname: "/" }} />
      )
    }

    let title = state.group_id === '' ?
    <h1>{i18n.t('group.new')}</h1> : <h1>{i18n.t('group.edit', {name: state.name})}</h1>;

    return (
      <div className='container'>
        <main className='column main'>
          <div className={style.form}>
            {title}
            <form onSubmit={this.handleSubmit}>
              <div>
                <input type='text' name='name' pattern='.{3,}' required value={state.name} autoComplete='off' placeholder='Name *' onChange={this.handleChange} />
              </div>
              <div className={style.body}>
                <CodeMirror
                  className='editor'
                  value={state.description}
                  options={{
                    mode: 'markdown',
                    theme: 'xq-light',
                    lineNumbers: true,
                    lineWrapping: true,
                    placeholder: 'Description'
                  }}
                  onBeforeChange={(editor, data, value) => this.handleDescriptionChange(editor, data, value)}
                />
              </div>
              <div>
                <button type="submit" className='btn topic' disabled={state.submitting}>
                  &nbsp;{i18n.t('general.submit')}
                </button>
              </div>
            </form>
          </div>
        </main>
        <aside className='column aside'>
        </aside>
      </div>
    )
  }
}

export default New;
