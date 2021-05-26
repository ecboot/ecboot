package com.github.ecboot.support;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.github.ecboot.enums.ResultEnum;
import lombok.Data;
import org.springframework.validation.BindingResult;

import java.util.Objects;

@Data
@JsonInclude(JsonInclude.Include.NON_NULL)
public class JsonResponse<T> {

    private String status;

    private T data;

    private Error error;

    private JsonResponse(T data) {
        this.status = "success";
        this.data = data;
    }

    private JsonResponse(Integer code, String message) {
        this.status = "failed";
        this.error = new Error(code, message);
    }

    public static <T> JsonResponse<T> succeed(T data) {
        return new JsonResponse<>(data);
    }

    public static <T> JsonResponse<T> failed(ResultEnum resultEnum) {
        return new JsonResponse<>(resultEnum.getCode(), resultEnum.getMessage());
    }

    public static <T> JsonResponse<T> failed(ResultEnum resultEnum, BindingResult bindingResult) {
        return new JsonResponse<>(resultEnum.getCode(),
                Objects.requireNonNull(bindingResult.getFieldError()).getField() + " " + bindingResult.getFieldError().getDefaultMessage());
    }

    @Data
    private static class Error {
        private Integer code;

        private String message;

        public Error(Integer code, String message) {
            this.code = code;
            this.message = message;
        }
    }
}