import style from './new.module.scss';
import React, { Component } from 'react';
import { Redirect } from 'react-router-dom';
import {UnControlled as CodeMirror} from 'react-codemirror2'
import showdown from 'showdown';
import showdownHighlight from 'showdown-highlight';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import Base64 from '../components/base64.js';
import API from '../api/index.js';
import Loading from '../components/loading.js';
import Button from '../components/button.js';
require('codemirror/lib/codemirror.css');
require('codemirror/theme/xq-light.css');
require('codemirror/mode/markdown/markdown.js');
const validate = require('uuid-validate');

class New extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.base64 = new Base64();
    this.converter = new showdown.Converter({ extensions: ['header-anchors', showdownHighlight] });
    this.instance = null;
    let categories = [];
    let d = window.localStorage.getItem('categories');
    if (!!d) {
      categories = JSON.parse(this.base64.decode(d));
    }
    let id = this.props.match.params.id;
    // false , 0 , '' , null , undefined , and NaN
    if (!id) {
      id = ''
    }
    this.state = {
      editor: 'deditor',
      topic_id: id,
      title: '',
      body: '',
      topic_type: 'POST',
      draft: false,
      categories: categories,
      preview: false,
      loading: true,
      submitting: false
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleCategoryClick = this.handleCategoryClick.bind(this);
    this.handleBodyChange = this.handleBodyChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleDraft = this.handleDraft.bind(this);
    this.handlePreview = this.handlePreview.bind(this);
    this.handleAction = this.handleAction.bind(this);
  }

  componentDidMount() {
    let id = 'draft'
    if (validate(this.state.topic_id)) {
      id = this.state.topic_id;
    }
    this.api.topic.show(id).then((resp) => {
      if (resp.error) {
        return;
      }
      let data = resp.data;
      if (!data) {
        data = {body: ''};
      }
      let l = data.body.split('\n').length;
      if (l > 13) {
        data.body += '\n'.repeat(3);
        data.editor = 'editor';
      }
      data.loading = false;
      this.setState(data);
    });
    this.api.category.index().then((resp) => {
      if (resp.error) {
        return;
      }
      const data = resp.data;
      let category_id = this.state.category_id;
      if (!category_id && data.length > 0) {
        category_id = data[0].category_id;
      }
      this.setState({categories: data, category_id: category_id});
    });
  }

  componentDidUpdate(prevProps, prevState) {
    if (this.props.location.pathname !== prevProps.location.pathname) {
      if (this.props.location.pathname === '/topics/new') {
        this.setState({
          topic_id: '',
          title: '',
          body: '',
          preview: false,
        });
      }
    }
  }

  handleChange(e) {
    const name = e.target.name;
    const value = e.target.type === 'checkbox' ? (e.target.checked ? 'LINK' : 'POST') : e.target.value;
    this.setState({
      [name]: value
    });
  }

  handleCategoryClick(e, value) {
    e.preventDefault();
    this.setState({category_id: value});
  }

  handleBodyChange(editor, data, value) {
    if (!data.origin) {
      return;
    }
    let cursor = editor.getCursor();
    let l = value.split('\n').length;
    let style = 'deditor';
    if (l > 13 || value.length > 1200) {
      if (value.charAt(value.length - 1) !== "\n") {
        value += '\n';
      }
      style = 'editor'
    }
    this.setState({body: value, editor: style}, () => {
      let indent = 0;
      if (data.text.length===1) {
        switch (data.text[0]) {
          case '**':
          case '****':
          case '~~~~':
            indent = data.text[0].length/2;
            break;
          default:
        }
      }
      this.instance.setCursor({line: cursor.line, ch: cursor.ch-indent});
    });
  }

  handlePreview(e) {
    e.preventDefault();
    this.setState({body_html: this.converter.makeHtml(this.state.body), preview: !this.state.preview});
  }

  handleAction(action, identity) {
    if (this.instance !== null) {
      let editor = this.instance;
      switch (action) {
        case 'heading':
          let cursor = editor.getCursor();
          editor.replaceRange('## ', {line: cursor.line, ch: 0}, {line: cursor.line, ch: 0}, '+input');
          break;
        case 'bold':
        case 'italic':
        case 'strikethrough':
          let t = editor.getSelection().trim();
          editor.replaceSelection(editor.getSelection().replace(t, `${identity}${t}${identity}`));
          break;
        case 'ol':
        case 'ul':
        case 'quote':
          let selections = editor.listSelections();
          selections.forEach(function(selection) {
            let pos = [selection.head.line, selection.anchor.line];
            if (selection.head.line > selection.anchor.line)  {
              pos = [selection.anchor.line, selection.head.line];
            }
            for (let i=pos[0]; i<=pos[1]; i++) {
              editor.replaceRange(identity, { line: i, ch: 0 }, { line: i, ch: 0 }, '+input');
            }
          });
          break;
        default:
      }
      this.instance.focus();
    }
  }

  handleSubmit(e) {
    e.preventDefault();
    if (this.state.submitting) {
      return;
    }
    this.setState({submitting: true, draft: false}, () => {
      this.submitForm();
    });
  }

  handleDraft(e) {
    e.preventDefault();
    this.setState({submitting: true, draft: true}, () => {
      this.submitForm();
    });
  }

  submitForm() {
    const history = this.props.history;
    if (validate(this.state.topic_id)) {
      this.api.topic.update(this.state.topic_id, this.state).then((resp) => {
        this.setState({submitting: false});
        if (resp.error) {
          return;
        }
        if (this.state.draft) {
          return;
        }
        history.push(`/topics/${resp.data.topic_id}`);
      });
      return;
    }
    this.api.topic.create(this.state).then((resp) => {
      if (resp.error) {
        return
      }
      this.setState({submitting: false});
      history.push(`/topics/${resp.data.topic_id}`);
    });
  }

  render() {
    const i18n = window.i18n;
    if (!this.api.user.loggedIn()) {
      return (
        <Redirect to={{ pathname: '/' }} />
      )
    }

    let state = this.state;
    const categories = state.categories.map((c) => {
      return (
        <span key={c.category_id} className={`${style.category} ${c.category_id === state.category_id ? style.active : ''}`} onClick={(e) => this.handleCategoryClick(e, c.category_id)}>{c.alias}</span>
      )
    });

    let title = <h1>{i18n.t('topic.title.new')}</h1>;
    if (validate(state.topic_id)) {
      title = <h1>{i18n.t('topic.title.edit', {name: state.title})}</h1>
    }

    const loadingView = (
      <div className={style.loading}>
        <Loading class='medium'/>
      </div>
    )

    const toolbar = TOOLBAR.map((data) => {
      return <FontAwesomeIcon className={style.action} icon={['fas', data.icon]} onClick={this.handleAction.bind(this, data.action, data.identity)} />
    });

    let form = (
      <form onSubmit={this.handleSubmit}>
        <div className={style.categories}>
          {categories}
        </div>
        <div>
          <input type='text' name='title' pattern='.{3,}' required value={state.title} autoComplete='off' placeholder={i18n.t('topic.placeholder.title')} onChange={this.handleChange} />
        </div>
        {
          state.topic_type === 'POST' &&
          <div className={style.actions}>
            <div className={style.toolbar}>
              { toolbar }
            </div>
            { state.preview && <FontAwesomeIcon className={style.eye} icon={['far', 'eye-slash']} onClick={this.handlePreview} /> }
            { !state.preview && <FontAwesomeIcon className={style.eye} icon={['far', 'eye']} onClick={this.handlePreview} /> }
            <a className={style.markdown} href='https://guides.github.com/features/mastering-markdown/' target='_blank' rel='noopener noreferrer'>
              <FontAwesomeIcon className={style.eye} icon={['fab', 'markdown']} />
            </a>
          </div>
        }
        {
          state.topic_type === 'POST' &&
          <div className={style.body}>
            {
              !state.preview &&
              <CodeMirror
                className={state.editor}
                value={state.body}
                options={{
                  mode: 'markdown',
                  theme: 'xq-light',
                  lineNumbers: true,
                  lineWrapping: true,
                  placeholder: 'Text (optional)'
                }}
                onChange={(editor, data, value) => this.handleBodyChange(editor, data, value)}
                editorDidMount={editor => {
                  this.instance = editor;
                  this.instance.refresh();
                }}
              />
            }
            {
              state.preview &&
              <article className={`md ${style.preview}`} dangerouslySetInnerHTML={{__html: state.body_html}}>
              </article>
            }
          </div>
        }
        <div className={style.upload}>
          <a href='https://imgur.com/upload' target='_blank' rel='noopener noreferrer'>Does not support upload image, please use imgur first.</a>
        </div>
        {
          state.topic_type === 'LINK' &&
          <div>
            <textarea name='body' rows='2' value={state.body} onChange={this.handleChange} className={style.link} placeholder={i18n.t('topic.placeholder.url')}></textarea>
          </div>
        }
        <div className={style.submit}>
          <Button type='submit' classes='submit' disabled={state.submitting} text={i18n.t('general.submit')} />
          {!state.submitting && (state.topic_id === "" || state.draft) && <span className={style.draft} onClick={this.handleDraft}>{i18n.t('general.draft')}</span>}
        </div>
      </form>
    )

    return (
      <div className='container'>
        <main className='column main'>
          {state.loading && loadingView}
          <div className={style.form}>
            {!state.loading && title}
            {!state.loading && form}
          </div>
        </main>
        <aside className='column aside'>
          <div className={style.title}>Rules</div>
          <ul className={style.rules} dangerouslySetInnerHTML={{__html: i18n.t('topic.rules')}}></ul>
        </aside>
      </div>
    )
  }
}

const TOOLBAR = [
  {icon: 'heading', action: 'heading', identity: ''},
  {icon: 'bold', action: 'bold', identity: '**'},
  {icon: 'italic', action: 'italic', identity: '*'},
  {icon: 'strikethrough', action: 'strikethrough', identity: '~~'},
  {icon: 'quote-left', action: 'quote', identity: '> '},
  {icon: 'list-ol', action: 'ol', identity: '1. '},
  {icon: 'list-ul', action: 'ul', identity: '* '},
]


export default New;
