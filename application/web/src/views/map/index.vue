<template>
  <div class="map">
    <el-input
      v-model="input"
      placeholder="请输入溯源码查询"
      style="width: 300px; margin-right: 15px"
      @keyup.enter="fetchRoute"
    />
    <el-button @click="fetchRoute">查询</el-button>
    <div id="container" style="width: 100%; height: 90vh" />
  </div>
</template>

<script>
import axios from "axios";
import { getFruitInfo } from "@/api/trace";

export default {
  data() {
    return {
      input: "",
      firstArr: [108.983569, 34.285675],
      map: null,
      marker: null,
      polyline: null,
      locationMapping: {},
      tracedata: null,
    };
  },
  mounted() {
    setTimeout(() => {
      this.initMap();
    }, 1000);
    this.loadLocationMapping();
  },
  methods: {
    initMap() {
      this.map = new AMap.Map("container", {
        resizeEnable: true,
        center: this.firstArr,
        zoom: 5,
      });
      this.marker = new AMap.Marker({
        map: this.map,
        position: this.firstArr,
        icon: "https://webapi.amap.com/images/car.png",
        offset: new AMap.Pixel(-26, -13),
        autoRotation: true,
        angle: -90,
      });
    },
    async loadLocationMapping() {
      try {
        this.locationMapping = require("@/location_mapping.json");
      } catch (error) {
        console.error("Error loading location mapping:", error);
        this.$message.error("加载地点映射数据失败");
      }
    },
    async fetchRoute() {
      if (!this.input) {
        this.$message.error("请输入有效的溯源码");
        return;
      }

      try {
        const formData = new FormData();
        formData.append("traceability_code", this.input);
        const res = await getFruitInfo(formData);

        if (res.code === 200) {
          this.tracedata = [];
          this.tracedata = JSON.parse(res.data);
          const locations = [];
          if (this.tracedata.farmer_input.fa_origin) {
            locations.push(this.tracedata.farmer_input.fa_origin);
          }
          if (
            this.tracedata.farmer_input.fa_origin &&
            this.tracedata.factory_input.fac_factoryName
          ) {
            locations.push(this.tracedata.factory_input.fac_factoryName);
          }
          if (
            this.tracedata.farmer_input.fa_origin &&
            this.tracedata.factory_input.fac_factoryName &&
            this.tracedata.shop_input.sh_shopAddress
          ) {
            locations.push(this.tracedata.shop_input.sh_shopAddress);
          }

          const coordinates = this.convertLocationsToCoordinates(locations);
          if (coordinates.length > 0) {
            this.updateMap(coordinates);
          } else {
            this.$message.error("未找到有效的坐标数据");
          }
        } else {
          this.$message.error(res.message);
        }
      } catch (error) {
        console.error("Error fetching route data:", error);
        this.$message.error("获取路线数据失败");
      }
    },
    convertLocationsToCoordinates(locations) {
      return locations
        .map((location) => {
          const coords = this.locationMapping[location];
          if (coords) {
            return coords;
          } else {
            console.warn(`No coordinates found for location: ${location}`);
            return null;
          }
        })
        .filter((coord) => coord !== null);
    },
    updateMap(coordinates) {
      if (this.polyline) {
        this.polyline.setMap(null);
      }
      if (this.marker) {
        this.marker.setMap(null);
      }

      // 创建折线，并设置路径和样式
      this.polyline = new AMap.Polyline({
        map: this.map,
        path: coordinates,
        showDir: true,
        strokeColor: "#77DDFF",
        strokeWeight: 6,
        lineJoin: "round",
      });

      // 获取路线的终点坐标
      const endPoint = coordinates[coordinates.length - 1];

      // 在路线终点处创建汽车图标，并设置位置和样式
      this.marker = new AMap.Marker({
        map: this.map,
        position: endPoint,
        icon: "https://webapi.amap.com/images/car.png",
        offset: new AMap.Pixel(-26, -13),
        autoRotation: true,
        angle: -90,
      });

      // 调整地图视野，使得整条路线都能显示在地图中
      this.map.setFitView();
    },
  },
};
</script>

<style lang="scss" scoped></style>
