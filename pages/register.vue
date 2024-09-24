<template>
  <div class="min-h-screen flex items-center justify-center">
    <Card class="w-1/2 min-w-96 shadow-md rounded-md p-6 dark">
      <CardHeader>
        <CardTitle class="text-2xl font-semibold text-center">
          User Registration
        </CardTitle>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="onSubmit">
          <div class="mb-4">
            <FormField v-slot="{ componentField }" name="username">
              <FormItem>
                <FormLabel for="username" class="text-right">
                  Username
                </FormLabel>
                <FormDescription>
                  This is your public display name
                </FormDescription>
                <FormControl>
                  <Input
                    v-bind="componentField"
                    id="username"
                    type="text"
                    placeholder="Enter your username"
                    class="col-span-3 dark"
                    @input="uniqueUsername = false"
                  />
                </FormControl>
                <ErrorMessage name="username" class="text-sm text-red-500" />
                <p v-if="uniqueUsername" class="text-red-500 text-sm mt-1">
                  Username is already in use
                </p>
              </FormItem>
            </FormField>

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
                    @input="uniqueEmail = false"
                  />
                </FormControl>
                <ErrorMessage name="email" class="text-sm text-red-500" />
                <p v-if="uniqueEmail" class="text-red-500 text-sm mt-1">
                  Email is already in use
                </p>
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
                  />
                </FormControl>
                <ErrorMessage name="password" class="text-sm text-red-500" />
              </FormItem>
            </FormField>

            <FormField v-slot="{ componentField }" name="confirmPassword">
              <FormItem>
                <FormLabel for="confirmPassword" class="text-right">
                  Confirm Password
                </FormLabel>
                <FormControl>
                  <Input
                    v-bind="componentField"
                    id="confirmPassword"
                    type="password"
                    placeholder="Confirm your password"
                    class="col-span-3 dark"
                  />
                </FormControl>
                <ErrorMessage
                  name="confirmPassword"
                  class="text-sm text-red-500"
                />
              </FormItem>
            </FormField>

            <Button type="submit" class="w-full mt-4">Register</Button>
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

const uniqueEmail = ref(false);
const uniqueUsername = ref(false);

// Define the validation schema using Yup
const schema = yup.object({
  username: yup
    .string()
    .required("Username is required")
    .min(3, "Username must be at least 3 characters"),
  email: yup
    .string()
    .required("Email is required")
    .email("Must be a valid email"),
  password: yup
    .string()
    .required("Password is required")
    .min(6, "Password must be at least 6 characters"),
  confirmPassword: yup
    .string()
    .required("Please confirm your password")
    .oneOf([yup.ref("password")], "Passwords must match"),
});

// Initialize useForm with the validation schema
const form = useForm({
  validationSchema: schema,
});

// Handle form submission
const onSubmit = form.handleSubmit(async (values) => {
  const registrationData = {
    username: values.username,
    email: values.email,
    password: values.password,
  };

  const { error } = await useAsyncData("register", () =>
    $fetch("http://localhost:8080/user/register", {
      method: "POST",
      body: registrationData,
    })
  );
  if (error.value && error.value.statusCode !== 409) {
    console.log(error.value);
    console.error("Registration error:", error.value);
    toast({
      title: "Registration failed!",
      description: "Please try again.",
    });
  } else if (error.value && error.value.statusCode === 409) {
    console.log("error.value.message", error.value.data);
    const cleanMessage = (error.value.data as string).trim();
    if (cleanMessage === "email is already in use") {
      uniqueEmail.value = true;
    } else if (cleanMessage === "username is already in use") {
      uniqueUsername.value = true;
    }
  } else {
    toast({
      title: "Registration successful!",
      description: "You have successfully registered.",
    });
    navigateTo("/");
  }
});
</script>
