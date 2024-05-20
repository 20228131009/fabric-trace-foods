<template>
  <div class="trace-container">
    <!-- 种植户 -->
    <div v-show="userType == '种植户'" class="table-container">
      <el-table :data="tracedata" border style="width: auto">
        <!-- 食品名称 -->
        <el-table-column
          prop="farmer_input.fa_fruitName"
          label="食品名称"
          width="auto"
        />
        <!-- 产地 -->
        <el-table-column
          prop="farmer_input.fa_origin"
          label="产地"
          width="auto"
        />
        <!-- 种植户名称 -->
        <el-table-column
          prop="farmer_input.fa_farmerName"
          label="种植户名称"
          width="auto"
        />
        <!-- 操作 -->
        <el-table-column label="操作">
          <template slot-scope="scope">
            <el-button size="mini">查看</el-button>
            <el-button type="danger" size="mini">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <!-- 工厂 -->
    <div v-show="userType == '工厂'" class="table-container">
      <el-table :data="tracedata" border style="width: auto">
        <!-- 商品名称 -->
        <el-table-column
          prop="factory_input.fac_productName"
          label="商品名称"
        />
        <!-- 生产批次 -->
        <el-table-column
          prop="factory_input.fac_productionbatch"
          label="生产批次"
        />
        <!-- 工厂名称 -->
        <el-table-column
          prop="factory_input.fac_factoryName"
          label="工厂名称"
        />
        <!-- 操作 -->
        <el-table-column label="操作">
          <template slot-scope="scope">
            <el-button size="mini">查看</el-button>
            <el-button type="danger" size="mini">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <!-- 运输司机 -->
    <div v-show="userType == '运输司机'" class="table-container">
      <el-table :data="tracedata" border style="width: auto">
        <!-- 司机姓名 -->
        <el-table-column prop="driver_input.dr_name" label="司机姓名" />
        <!-- 联系电话 -->
        <el-table-column prop="driver_input.dr_phone" label="联系电话" />
        <!-- 运输记录 -->
        <el-table-column prop="driver_input.dr_transport" label="运输记录" />
        <!-- 操作 -->
        <el-table-column label="操作">
          <template slot-scope="scope">
            <el-button size="mini">查看</el-button>
            <el-button type="danger" size="mini">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <!-- 商店 -->
    <div v-show="userType == '商店'" class="table-container">
      <el-table :data="tracedata" border style="width: auto">
        <!-- 入库时间 -->
        <el-table-column prop="shop_input.sh_storeTime" label="入库时间" />
        <!-- 商店名称 -->
        <el-table-column prop="shop_input.sh_shopName" label="商店名称" />
        <!-- 商店电话 -->
        <el-table-column prop="shop_input.sh_shopPhone" label="商店电话" />
        <!-- 操作 -->
        <el-table-column label="操作">
          <template slot-scope="scope">
            <el-button size="mini">查看</el-button>
            <el-button type="danger" size="mini">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <!-- 消费者 -->
    <div v-show="userType == '消费者'" class="input-area">无法进行上链哦</div>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import {
  getFruitInfo,
  getFruitList,
  getAllFruitInfo,
  getFruitHistory,
} from "@/api/trace";

export default {
  name: "Trace",
  data() {
    return {
      tracedata: [],
      loading: false,
      input: "",
      imageSrc: "",
      // traceability_code: this.traceabilityCode(),
      // // 使用函数调用获取 imagePath
      // imagePath:
      //   "http://localhost:9090/dist/upload/" + this.traceabilityCode() + ".jpg",
    };
  },
  computed: {
    ...mapGetters(["name", "userType"]),
    ...mapGetters(["selectedPhotos"]),
    // traceabilityCode() {
    //   // 假设 traceability_code 是从 tracedata 中的某个字段获取的
    //   // 如果这里不适用，请替换成合适的字段路径
    //   return this.tracedata.traceability_code; // 由于是数组，这里使用了join来合并
    // },
  },
  created() {
    getFruitList().then((res) => {
      this.tracedata = JSON.parse(res.data);
    });
  },
  methods: {
    AllFruitInfo() {
      getAllFruitInfo().then((res) => {
        this.tracedata = JSON.parse(res.data);
      });
    },
    FruitHistory() {
      getFruitHistory().then((res) => {
        // console.log(res)
      });
    },
    FruitInfo() {
      var formData = new FormData();
      formData.append("traceability_code", this.input);
      getFruitInfo(formData).then((res) => {
        if (res.code === 200) {
          // console.log(res)
          this.tracedata = [];
          this.tracedata[0] = JSON.parse(res.data);
          return;
        } else {
          this.$message.error(res.message);
        }
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.trace-container {
  margin: 30px;
}

.table-container {
  margin-bottom: 20px;
}

.el-table {
  width: 100%;
}
</style>
