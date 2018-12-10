<template lang="pug">
b-col(md='12')
  b-card
    div(slot='header')
      strong Device Details

    b-form
      b-form-group(description='IP address or domain name', label='Host', :label-cols='3', :horizontal='true')
        b-form-input#hostname(type='text', v-model="cfg.host")
      b-form-group(description='TCP port for incoming gRPC.', label='Port', :label-cols='3', :horizontal='true')
        b-form-input#port(type='number', v-model="cfg.port")
      b-form-group(description='gRPC username', label='Username', :label-cols='3', :horizontal='true')
        b-form-input#username(type='text', autocomplete='username' v-model="cfg.user")
      b-form-group(description='Password', label='Password', :label-cols='3', :horizontal='true')
        b-form-input#password(type='password', v-model="cfg.password")
      b-form-group(description='Unique string to identify this particular client', label='Client ID', :label-cols='3', :horizontal='true')
        b-form-input#cid(type='text', autocomplete='clientID' v-model="cfg.cid")
      b-form-group(description='HTTP2 window size', label='Window Size', :label-cols='3', :horizontal='true')
        b-form-input#ws(type='number', v-model="cfg.ws")
      b-form-group(description='How often do we recive data', label='Frequency', :label-cols='3', :horizontal='true')
        b-form-input#fq(type='number', v-model="cfg.freq")  
      
      div(slot='footer')
        b-button(size='sm', , v-on:click="add" variant='primary')
          i.fa.fa-dot-circle-o
          |  Submit
        b-button(size='sm', v-on:click="cancel", variant='danger')
          i.fa.fa-ban
          |  Cancel

</template>

<script>
import axios from 'axios'
import router from '@/router'

export default {
  name: 'newdevice',
  props: {
    item: {
      type: Object,
      requred: false
    }  
  },
  data () {    
    if (this.item == null){
      return {      
        errors: [],
        cfg: {
          host:     '',
          port:     50051,
          user:     '',
          password: '',
          cid:      '',
          ws:       524288,
          uuid: '',
          freq: 2000			  
        }   
      }
    } else {     
      return {
        errors: [],
        cfg: {
          host: this.item.host,
			    port: this.item.port,
			    user: this.item.user,
			    password: this.item.password,
			    cid: this.item.cid,
          ws: this.item.ws,
          freq: this.item.paths[0].freq,
          uuid: this.item.uuid
        }
      }
    }  
  },
  created () {
    
  },
  methods: {   
    cancel () {
      router.push('/devices')
    },
    add () {
      axios
        .post('http://' + window.location.hostname + ':8888/v1/device',        
        {
          host: this.cfg.host,
			    port: parseInt(this.cfg.port, 10),
			    user: this.cfg.user,
			    password: this.cfg.password,
			    cid: this.cfg.cid,
          ws: parseInt(this.cfg.ws, 10),
          uuid: this.cfg.uuid,
          paths: [
          {
            path: "/interfaces",
					  freq: parseInt(this.cfg.freq, 10)
				  }],
			  })           
        .then(response => {   
          router.push('/devices')     
          console.log(response)        
        })
        .catch(err => {
          console.log(err)
        })
    }                      
  } 
}
</script>
