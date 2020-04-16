


def lazy_pick(coord_list, tag_data, region_data, num_slice, intensity_param):
    # Height of slice
    segment_size_y = 1.0 / num_slice

    for idx, point in enumerate(coord_list, start=0):
        normalized_counter = min(
            max(tag_data[idx][0] / intensity_param, 0.0), 1.0)
        if(normalized_counter == 0.0):
            continue
        # Clamped (0, 255)
        effectpow = int(normalized_counter * intensity_param * 255)
        region_data[idx][1] = effectpow

        for iter in range(num_slice):
            if(point[1] > (iter * segment_size_y) and point[1] < ((1 + iter) * segment_size_y)):
                region_data[idx][0] = iter
