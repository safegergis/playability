<template>
  <div>
    <Dialog>
      <DialogTrigger as-child>
        <Button variant="outline" class="dark">
          Submit Accessibility Report
        </Button>
      </DialogTrigger>
      <DialogContent
        class="bg-primary text-primary-foreground border-primary sm:max-w-[625px]"
      >
        <DialogHeader class="">
          <DialogTitle>Submit Accessibility Report</DialogTitle>
          <DialogDescription>
            Please provide information about the game's accessibility features.
          </DialogDescription>
        </DialogHeader>
        <form @submit.prevent="onSubmit">
          <div class="grid gap-4 py-4">
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="colorBlind" class="text-right"
                >Color Blind Support</Label
              >
              <Select
                id="colorBlind"
                v-model="formData.colorBlind"
                class="col-span-3"
              >
                <SelectTrigger class="dark">
                  <SelectValue placeholder="Select option" />
                </SelectTrigger>
                <SelectContent class="dark">
                  <SelectItem value="false">There is no support</SelectItem>
                  <SelectItem value="unknown">I don't know</SelectItem>
                  <SelectItem value="limited"
                    >There is limited support</SelectItem
                  >
                  <SelectItem value="true">There is full support</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="closedCaptions" class="text-right"
                >Closed Caption Support</Label
              >
              <Select
                id="closedCaptions"
                v-model="formData.closedCaptions"
                class="col-span-3"
              >
                <SelectTrigger class="dark">
                  <SelectValue placeholder="Select option" />
                </SelectTrigger>
                <SelectContent class="dark">
                  <SelectItem value="false">There is no support</SelectItem>
                  <SelectItem value="unknown">I don't know</SelectItem>
                  <SelectItem value="limited"
                    >There is limited support</SelectItem
                  >
                  <SelectItem value="true">There is full support</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="controllerSupport" class="text-right"
                >Full Controller Support</Label
              >
              <Select
                id="controllerSupport"
                v-model="formData.controllerSupport"
                class="col-span-3"
              >
                <SelectTrigger class="dark">
                  <SelectValue placeholder="Select option" />
                </SelectTrigger>
                <SelectContent class="dark">
                  <SelectItem value="false">There is no support</SelectItem>
                  <SelectItem value="unknown">I don't know</SelectItem>
                  <SelectItem value="limited"
                    >There is limited support</SelectItem
                  >
                  <SelectItem value="true">There is full support</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="controllerRemapping" class="text-right"
                >Controller Remapping</Label
              >
              <Select
                id="controllerRemapping"
                v-model="formData.controllerRemapping"
                class="col-span-3"
              >
                <SelectTrigger class="dark">
                  <SelectValue placeholder="Select option" />
                </SelectTrigger>
                <SelectContent class="dark">
                  <SelectItem value="false">There is no support</SelectItem>
                  <SelectItem value="unknown">I don't know</SelectItem>
                  <SelectItem value="limited"
                    >There is limited support</SelectItem
                  >
                  <SelectItem value="true">There is full support</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="score" class="text-right"> Score (out of 10) </Label>
              <Input
                id="score"
                v-model="formData.score"
                type="number"
                min="0"
                max="10"
                class="col-span-3 dark"
              />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="report" class="text-right">Additional Comments</Label>
              <Textarea
                id="report"
                v-model="formData.report"
                placeholder="Optional report for others"
                class="col-span-3 dark"
              />
            </div>
          </div>
          <DialogFooter>
            <Button type="submit">Submit Report</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script lang="ts" setup>
const props = defineProps<{
  gameId: number;
}>();

const formData = ref({
  colorBlind: "",
  closedCaptions: "",
  controllerSupport: "",
  controllerRemapping: "",
  score: 0,
  report: "",
});
const ReportRow = computed(() => {
  return {
    game_id: props.gameId,
    closed_captions: formData.value.closedCaptions,
    color_blind: formData.value.colorBlind,
    full_controller_support: formData.value.controllerRemapping,
    controller_remapping: formData.value.controllerSupport,
    score: formData.value.score,
    report: formData.value.report,
  };
});

const jwt = useCookie("jwt");
console.log(jwt.value);
const onSubmit = async () => {
  console.log(formData.value);
  await useFetch("/api/report/report", {
    method: "POST",
    body: ReportRow.value,
  });
};
</script>

<style scoped>
/* Add any component-specific styles here */
</style>
