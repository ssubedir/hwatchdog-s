<template>
  <div class="index">
    <div class="container-fluid">
      <!-- Button trigger modal -->

      <h1>CPUs Data</h1>
      <div class="float-right">
        <button
          type="button"
          class="btn btn-primary"
          data-toggle="modal"
          data-target="#exampleModalCenter"
        >
          <i class="fa fa-filter" aria-hidden="true"></i>
        </button>
      </div>
      <br />
      <br />
      <table class="table">
        <thead>
          <tr>
            <th scope="col">Name</th>
            <th scope="col">Price</th>
            <th scope="col">Last Checked</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(cpu, index) in hdata.Cpus">
            <td>{{ cpu.Name }}</td>
            <td>${{ cpu.Price == null ? "0.00" : cpu.Price }}</td>
            <td>
              {{
                dateformat(cpu.Time) == "2020 years ago" ||
                dateformat(cpu.Time) == "2021 years ago"
                  ? "n/a"
                  : dateformat(cpu.Time)
              }}
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Modal -->
      <div
        class="modal fade"
        id="exampleModalCenter"
        tabindex="-1"
        role="dialog"
        aria-labelledby="exampleModalCenterTitle"
        aria-hidden="true"
      >
        <div class="modal-dialog modal-dialog-centered" role="document">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title" id="exampleModalLongTitle">
                <i class="fa fa-filter" aria-hidden="true"></i>
                Filter
              </h5>
              <button
                type="button"
                class="close"
                data-dismiss="modal"
                aria-label="Close"
              >
                <span aria-hidden="true">&times;</span>
              </button>
            </div>
            <div class="modal-body">
              <form>
                <div class="form-group row">
                  <label for="exampleFormControlSelect1">SITE</label>
                  <select class="form-control" id="filtersite">
                    <option>amazon</option>
                  </select>
                </div>
                <div class="form-group row">
                  <label for="exampleFormControlSelect1">LOCATION</label>
                  <select class="form-control" id="filterlocation">
                    <option>us</option>
                    <option>ca</option>
                  </select>
                </div>
              </form>
            </div>
            <div class="modal-footer">
              <button
                type="button"
                class="btn btn-secondary"
                data-dismiss="modal"
                id="filterclose"
              >
                Close
              </button>
              <button type="button" class="btn btn-primary" v-on:click="updatefilter()">
                Save changes
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import moment from "moment";
export default {
  name: "home",
  data() {
    return {
      hdata: {},
      filter: {},
    };
  },
  methods: {
    fetchdata() {
        var filt = JSON.parse(localStorage.hwfilter);
        console.log(filt);
      fetch("http://localhost:9001/api/"+filt.site+"?location="+filt.location)
        .then((response) => response.json())
        .then((data) => {
          this.hdata = data;
          //console.log(data);
        });
    },
    currentDateTime() {
      return moment().format("MMMM Do YYYY, h:mm:ss a");
    },
    dateformat(d) {
      return moment(d).fromNow();
    },
    updatefilter() {
     
      var s = document.getElementById("filtersite").value;
      var l = document.getElementById("filterlocation").value;

        var f = {};

        f = {
          site: s,
          location: l,
        };

        localStorage.hwfilter = JSON.stringify(f);

        
      var fc = document.getElementById("filterclose").click();

    },
    getfilter(){
      var filt = JSON.parse(localStorage.hwfilter);
      return filt;
    }
  },
  created() {
    this.fetchdata();
  },
  mounted() {
    if (localStorage.hwfilter) {
      this.filter = localStorage.hwfilter;
    } else {
      var f = {};

      f = {
        site: "amazon",
        location: "us",
      };

      localStorage.hwfilter = JSON.stringify(f);
    }
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
