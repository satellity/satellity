require('codemirror/lib/codemirror.css');
require('codemirror/theme/xq-light.css');
require('codemirror/mode/markdown/markdown.js');
import 'react-image-crop/lib/ReactCrop.scss';
import style from './new.scss';
import React, { Component } from 'react';
import { Redirect } from 'react-router-dom';
import {Controlled as CodeMirror} from 'react-codemirror2'
import ReactCrop from 'react-image-crop';
import API from '../api/index.js';
const validate = require('uuid-validate');

class New extends Component {
  constructor(props) {
    super(props);
    let id = this.props.match.params.id;
    if (!id) {
      id = ''
    }
    this.state = {
      group_id: id,
      name: '',
      description: '',
      image: '',
      cover_url: '',
      submitting: false,
      loading: false
    }

    this.api = new API();
    this.handleClick = this.handleClick.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.handleDescriptionChange = this.handleDescriptionChange.bind(this);
    this.handleFileChange = this.handleFileChange.bind(this);
    this.onImageLoaded = this.onImageLoaded.bind(this);
    this.onCropComplete = this.onCropComplete.bind(this);
    this.onCropChange = this.onCropChange.bind(this);
    this.makeClientCrop = this.makeClientCrop.bind(this);
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

  handleClick(e) {
    this.refs.file.click();
  }

  handleChange(e) {
    const {target: {name, value}} = e;
    this.setState({
      [name]: value
    });
  }

  handleDescriptionChange(editor, data, value) {
    this.setState({description: value});
  }

  handleFileChange(e) {
    if (e.target.files && e.target.files.length > 0) {
      const reader = new FileReader();
      reader.addEventListener('load', () => {
        this.setState({ image: reader.result });
      });
      reader.readAsDataURL(e.target.files[0]);
    }
  }

  onImageLoaded(image) {
    this.imageRef = image;
  }

  onCropComplete(crop) {
    this.makeClientCrop(crop);
  };

  onCropChange(crop, percentCrop) {
  };

  makeClientCrop(crop) {
    if (this.imageRef) {
      crop.width = this.imageRef.naturalWidth;
      crop.height = this.imageRef.naturalHeight;
      if (crop.width < 600 || crop.height < 300) {
        return;
      }
      if (crop.width > crop.height * 2) {
        crop.width = crop.height * 2;
      } else {
        crop.height = crop.width / 2;
      }
      const canvas = document.createElement("canvas");
      canvas.width = crop.width;
      canvas.height = crop.height;
      const ctx = canvas.getContext("2d");
      ctx.drawImage(
        this.imageRef,
        0,
        0,
        crop.width,
        crop.height,
        0,
        0,
        crop.width,
        crop.height
      );
      this.setState({cover_url: canvas.toDataURL()});
    }
  }

  handleSubmit(e) {
    e.preventDefault();
    if (this.state.submitting) {
      return
    }
    this.setState({submitting: true}, () => {
      const history = this.props.history;
      let i = this.state.cover_url.indexOf(",");
      const data = {name: this.state.name, description: this.state.description, cover: this.state.cover_url.slice(i+1)};
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
                <input type='text' name='name' pattern='.{3,}' required value={state.name} autoComplete='off' placeholder='Group Name *' onChange={this.handleChange} />
              </div>
              <div className={style.cropContainer}>
                <input type='file' ref='file' className={style.file} onChange={this.handleFileChange} />
                <ReactCrop
                  src={state.image}
                  crop={{aspect: 2}}
                  className={style.crop}
                  onImageLoaded={this.onImageLoaded}
                  onComplete={this.onCropComplete}
                  onChange={this.onCropChange}
                />
                <div>
                  {
                    !state.cover_url && (
                      <div className={style.box} onClick={this.handleClick}>Cover</div>
                    )
                  }
                </div>
                <div className={style.image}>
                  {
                    state.cover_url && (
                      <img src={state.cover_url} onClick={this.handleClick} className={style.cover} />
                    )
                  }
                </div>
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
