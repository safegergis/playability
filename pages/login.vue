<template>
  <div class="min-h-screen flex items-center justify-center">
    <Card class="w-1/2 min-w-96 shadow-md rounded-md p-6 dark">
      <CardHeader>
        <CardTitle class="text-2xl font-semibold text-center">Login</CardTitle>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="onSubmit">
          <div class="mb-4">
            <FormField v-slot="{ componentField }" name="email">
              <FormItem>
                <FormLabel for="email" class="text-right">Email</FormLabel>
                <FormControl>
                  <Input
                    v-bind="componentField"
                    id="email"
                    type="email"
                    placeholder="Enter your email"
                    class="col-span-3 dark"
                    @input="InvalidLogin = false"
                  />
                </FormControl>
                <ErrorMessage name="email" class="text-sm text-red-500" />
              </FormItem>
            </FormField>

            <FormField v-slot="{ componentField }" name="password">
              <FormItem>
                <FormLabel for="password" class="text-right">
                  Password
                </FormLabel>
                <FormControl>
                  <Input
                    v-bind="componentField"
                    id="password"
                    type="password"
                    placeholder="Enter your password"
                    class="col-span-3 dark"
                    @input="InvalidLogin = false"
                  />
                </FormControl>
                <ErrorMessage name="password" class="text-sm text-red-500" />
              </FormItem>
            </FormField>

            <p v-if="InvalidLogin" class="text-red-500 text-sm mt-1">
              Invalid email or password
            </p>
            <Button type="submit" class="w-full mt-4">Login</Button>
          </div>
        </form>
      </CardContent>
    </Card>
  </div>
</template>

<script lang="ts" setup>
import { ErrorMessage, useForm } from "vee-validate";
import * as yup from "yup";
import { useToast } from "@/components/ui/toast";

const { toast } = useToast();

const InvalidLogin = ref(false);

// Define the validation schema using Yup
const schema = yup.object({
  email: yup
    .string()
    .required("Email is required")
    .email("Must be a valid email"),
  password: yup
    .string()
    .required("Password is required")
    .min(6, "Password must be at least 6 characters"),
});

// Initialize useForm with the validation schema
const form = useForm({
  validationSchema: schema,
});

// Handle form submission
const onSubmit = form.handleSubmit(async (values) => {
  const { data, error } = await useAsyncData("login", () =>
    $fetch("/api/auth/login", {
      method: "POST",
      body: values,
    })
  );
  // Check if the error is a 401
  if (error.value && error.value.statusCode !== 401) {
    //if not 401 then show error
    console.log(error.value);
    console.error("Login error:", error.value);
    toast({
      title: "Login failed!",
      description: "Please try again.",
    });
  } else if (error.value && error.value.statusCode === 401) {
    //if 401 then show invalid login
    console.log("error.value.message", error.value.data);
    InvalidLogin.value = true;
  } else {
    //success
    console.log("Login successful! ", data.value);
    toast({
      title: "Login successful!",
      description: "You have successfully logged in.",
    });
    navigateTo("/");
  }
});
</script>
