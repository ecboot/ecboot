package com.github.ecboot.service.activity;

import com.github.ecboot.entity.FavourableActivity;
import com.github.ecboot.repository.FavourableActivityRepository;
import com.github.ecboot.repository.GoodsActivityRepository;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;


public interface FavourableActivityService {

    GoodsActivityRepository goodsActivityRepository;

    FavourableActivityRepository favourableActivityRepository;

    public Page<FavourableActivity> findAll(Pageable pageable) {
        return null;
    }
}
