<template>
  <div>
    <div>
      <h1>Dashboard</h1>
    </div>

    <div>
      <div class="row">
        <div class="col-6">
          <table class="table table-bordered">
            <!-- <thead>
              <tr>
                <th scope="col">ID</th>
                <th>Name</th>
                <th>Status</th>
              </tr>
            </thead> -->
            <tbody>
              <tr>
                <td>status</td>
                <td>healthy</td>
              </tr>
              <tr>
                <td>version</td>
                <td>0.1.0</td>
              </tr>
              <tr>
                <td>maintenance</td>
                <td>false</td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="col-6">
          <div class="row">
            <div class="col-2">
              <p>cpu</p>
            </div>
            <div class="col">
              <div
                class="progress"
                role="progressbar"
                aria-valuenow="20"
                aria-valuemin="0"
                aria-valuemax="100"
                style="height: 25px"
              >
                <div class="progress-bar bg-success" style="width: 20%">20%</div>
              </div>
            </div>
          </div>

          <div class="row">
            <div class="col-2">
              <p>memory</p>
            </div>
            <div class="col">
              <div
                class="progress"
                role="progressbar"
                aria-valuenow="75"
                aria-valuemin="0"
                aria-valuemax="100"
                style="height: 25px"
              >
                <div class="progress-bar bg-warning" style="width: 75%">75%</div>
              </div>
            </div>
          </div>

          <div class="row">
            <div class="col-2">
              <p>disk</p>
            </div>
            <div class="col">
              <div
                class="progress"
                role="progressbar"
                aria-valuenow="85"
                aria-valuemin="0"
                aria-valuemax="100"
                style="height: 25px"
              >
                <div class="progress-bar bg-danger" style="width: 85%">85%</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
