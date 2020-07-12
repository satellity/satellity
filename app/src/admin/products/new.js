import style from './new.module.scss';
import React, { Component } from 'react';
import { Redirect } from 'react-router-dom';
import Loading from '../../components/loading.js';
const validate = require('uuid-validate');

class New extends Component {
  constructor(props) {
    super(props);

    this.api = window.api;
    let id = this.props.match.params.id;
    this.state = {
      product_id: !id ? '' : id,
      name: '',
      cover: '',
      source: '',
      body: '',
      tags: [],
      tags_str: '',
      submitting: false,
      redirect: false,
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDidMount() {
    if (this.state.product_id !== '') {
      this.api.product.show(this.state.product_id).then((resp) => {
        if (resp.error) {
          return
        }
        let data = resp.data;
        data.tags_str = data.tags.join();
        data.cover = data.cover_url;
        this.setState(data);
      });
    }
  }

  handleChange(e) {
    const name = e.target.name;
    this.setState({
      [name]: e.target.value,
    });
  }

  handleSubmit(e) {
    e.preventDefault();
    if (this.state.submitting) {
      return;
    }
    const tags = this.state.tags_str.split(',').map((e) => {
      return e.trim();
    });
    if (validate(this.state.product_id)) {
      this.setState({submitting: true, tags: tags}, () => {
        this.api.product.admin.update(this.state).then((resp) => {
          if (resp.error) {
            return;
          }
          this.setState({submitting: false, redirect: true});
        });
      })
      return;
    }
    this.setState({submitting: true, tags: tags}, () => {
      this.api.product.admin.create(this.state).then((resp) => {
        if (resp.error) {
          return;
        }
        this.setState({submitting: false, redirect: true});
      });
    });
  }

  render() {
    let state = this.state;

    if (state.redirect) {
      return (
        <Redirect to='/admin/products' />
      )
    }

    return (
      <div>
        <h1 className='welcome'>
          New Product
        </h1>
        <div className='panel'>
          <div>
            { state.cover !== '' && <img className={style.cover} src={state.cover} alt={state.name} /> }
          </div>
          <form onSubmit={this.handleSubmit}>
            <div>
              <label name='name'>Name *</label>
              <input type='text' name='name' pattern='.{1,}' required value={state.name} autoComplete='off' onChange={this.handleChange} />
            </div>
            <div>
              <label name='cover'>Cover *</label>
              <input type='text' name='cover' required value={state.cover} autoComplete='off' onChange={this.handleChange} />
            </div>
            <div>
              <label name='source'>Source * { state.source !== '' && <a href={state.source} rel='noopener noreferrer' target='_blank'>Go</a>}</label>
              <input type='text' name='source' required value={state.source} autoComplete='off' onChange={this.handleChange} />
            </div>
            <div>
              <label name='body'>Body *</label>
              <textarea type='text' name='body' required value={state.body} onChange={this.handleChange}/>
            </div>
            <div>
              <label name='body'>tags *</label>
              <input type='text' name='tags_str' required value={state.tags_str} autoComplete='off' onChange={this.handleChange} />
            </div>
            <div className={style.submit}>
              <button className='btn submit blue' disabled={state.submitting}>
                { state.submitting && <Loading class='small white'/> }
                &nbsp;SUBMIT
              </button>
            </div>
          </form>
        </div>
      </div>
    )
  }
}

export default New;
