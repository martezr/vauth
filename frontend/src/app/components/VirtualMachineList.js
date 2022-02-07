import React from 'react';
import axios from 'axios';

class VirtualMachineList extends React.Component {
  state = {
    persons: []
  }

  componentDidMount() {
    axios.get(`https://jsonplaceholder.typicode.com/users`)
      .then(res => {
        const persons = res.data;
        this.setState({ persons });
      })
  }

  render() {
    return (
      <p>test</p>
    );
  }
}

export default VirtualMachineList;
