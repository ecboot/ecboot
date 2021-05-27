package com.github.ecboot.service.activity.impl;

import com.github.ecboot.repository.FavourableActivityRepository;
import com.github.ecboot.repository.GoodsActivityRepository;
import com.github.ecboot.service.activity.FavourableActivityService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class FavourableActivityServiceImpl implements FavourableActivityService {

    @Autowired
    FavourableActivityRepository favourableActivityRepository;

    @Autowired
    GoodsActivityRepository goodsActivityRepository;
}
