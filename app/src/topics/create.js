import style from './create.module.scss';
import React, { Component } from 'react';
import Button from '../components/button.js';

class Create extends Component {
  render() {
    const i18n = window.i18n;

    return (
      <div className={style.new}>
        <Button type='link' action='/topics/new' text={i18n.t('topic.new')} classes='button' />
      </div>
    )
  }
}

export default Create;
