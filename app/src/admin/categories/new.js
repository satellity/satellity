import React, {Component} from 'react';
import API from '../../api/index.js';
import Loading from '../../components/loading.js';

class AdminCategoryNew extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.state = {name: '', alias: '', position: '', description: '', submitting: false};
  }

  handleChange(e) {
    const target = e.target;
    const name = target.name;
    this.setState({
      [name]: target.value,
    });
  }

  handleSubmit(e) {
    e.preventDefault();
    this.setState({submitting: true});
    this.api.category.admin.create(this.state).then(() => {
      this.setState({submitting: false});
    });
  }

  render() {
    return (
      <div>
        <h1 className='welcome'>
          Create a new category
        </h1>
        <div className='panel'>
          <form onSubmit={(e) => handleSubmit(e)}>
            <div>
              <label name='name'>Name *</label>
              <input type='text' name='name' value={state.name} autoComplete='off' onChange={(e) => handleChange(e)}/>
            </div>
            <div>
              <label name='alias'>Alias *</label>
              <input type='text' name='alias' value={state.alias} autoComplete='off' onChange={(e) => handleChange(e)}/>
            </div>
            <div>
              <label name='position'>Position</label>
              <input type='number' name='position' value={state.position} autoComplete='off' onChange={(e) => handleChange(e)}/>
            </div>
            <div>
              <label name='name'>Description *</label>
              <textarea type='text' name='description' value={state.description} onChange={(e) => handleChange(e)} key='description'/>
            </div>
            <div className='action'>
              <button className='btn submit blue' disabled={state.submitting}>
                { state.submitting && <Loading class='small white'/> }
                &nbsp;SUBMIT
              </button>
            </div>
          </form>
        </div>
      </div>
    );
  }
}

export default AdminCategoryNew;
