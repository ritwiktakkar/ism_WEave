import math


def divide_into_lots(total, max_lots=5):
    if total <= max_lots:
        return [1] * total
    base_lot_size = total // max_lots
    remainder = total % max_lots
    lots = [base_lot_size] * max_lots
    for i in range(remainder):
        lots[i] += 1
    return lots


def calculate_supply_chain(order_size):
    # Constants and conversion factors
    COTTON_BALE_WEIGHT = 480  # lbs
    YARN_COUNT = 30
    YARN_CONE_WEIGHT = 5  # lbs
    UNFINISHED_FABRIC_LENGTH = 50  # yards
    UNFINISHED_FABRIC_WEIGHT = 16.74  # lbs
    UNFINISHED_FABRIC_WIDTH = 60  # inches
    FINISHED_FABRIC_LENGTH = 47.5  # yards
    FINISHED_FABRIC_WEIGHT = 15.90  # lbs
    FINISHED_FABRIC_WIDTH = 58.8  # inches
    BUTTON_WEIGHT = 0.00165  # lbs
    SHIRT_WEIGHT = 0.554  # lbs
    BUTTONS_PER_SHIRT = 7
    CARTON_CAPACITY = 20  # shirts
    CONTAINER_CAPACITY = 400  # cartons

    # Waste factors (adjusted to increase overall loss)
    COTTON_TO_YARN_YIELD = 0.85
    YARN_TO_UNFINISHED_FABRIC_YIELD = 0.90
    UNFINISHED_TO_FINISHED_FABRIC_YIELD = 0.95
    FINISHED_FABRIC_TO_CUT_PARTS_YIELD = 0.85

    # 1) Cotton bale calculation
    total_shirt_weight = order_size * SHIRT_WEIGHT
    cotton_needed = total_shirt_weight / (
        COTTON_TO_YARN_YIELD
        * YARN_TO_UNFINISHED_FABRIC_YIELD
        * UNFINISHED_TO_FINISHED_FABRIC_YIELD
        * FINISHED_FABRIC_TO_CUT_PARTS_YIELD
    )
    bales_needed = math.ceil(cotton_needed / COTTON_BALE_WEIGHT)
    actual_cotton_weight = bales_needed * COTTON_BALE_WEIGHT

    # 2) Cotton bale lots (divided into at most 5 lots)
    bale_lot_sizes = divide_into_lots(bales_needed)
    bale_lots = [size * COTTON_BALE_WEIGHT for size in bale_lot_sizes]

    # 3) Cotton yarn cones
    yarn_weight = actual_cotton_weight * COTTON_TO_YARN_YIELD
    cones_needed = math.ceil(yarn_weight / YARN_CONE_WEIGHT)

    # 4) Cotton yarn lots
    yarn_lot_sizes = divide_into_lots(cones_needed)
    yarn_lots = [size * YARN_CONE_WEIGHT for size in yarn_lot_sizes]

    # 5) Unfinished fabric pieces
    unfinished_fabric_weight = yarn_weight * YARN_TO_UNFINISHED_FABRIC_YIELD
    unfinished_fabric_pieces = math.ceil(
        unfinished_fabric_weight / UNFINISHED_FABRIC_WEIGHT
    )

    # 6) Unfinished fabric lots
    unfinished_fabric_lot_sizes = divide_into_lots(unfinished_fabric_pieces)
    unfinished_fabric_lots = [
        size * UNFINISHED_FABRIC_WEIGHT for size in unfinished_fabric_lot_sizes
    ]

    # 7) Finished fabric pieces
    finished_fabric_weight = (
        unfinished_fabric_weight * UNFINISHED_TO_FINISHED_FABRIC_YIELD
    )
    finished_fabric_pieces = math.ceil(finished_fabric_weight / FINISHED_FABRIC_WEIGHT)

    # 8) Finished fabric lots
    finished_fabric_lot_sizes = divide_into_lots(finished_fabric_pieces)
    finished_fabric_lots = [
        size * FINISHED_FABRIC_WEIGHT for size in finished_fabric_lot_sizes
    ]

    # 9) Cut parts
    cut_parts_weight = finished_fabric_weight * FINISHED_FABRIC_TO_CUT_PARTS_YIELD
    cut_parts = {
        "front_panel": math.ceil(order_size * 1.1),  # Assuming 10% extra for waste
        "back_panel": math.ceil(order_size * 1.1),
        "left_sleeve": math.ceil(order_size * 1.1),
        "right_sleeve": math.ceil(order_size * 1.1),
        "collar": math.ceil(order_size * 1.1),
        "front_pocket": math.ceil(order_size * 1.1),
    }
    cut_part_weight = cut_parts_weight / (sum(cut_parts.values()) / order_size)

    # 10) Buttons
    buttons_needed = order_size * BUTTONS_PER_SHIRT

    # 11) Shirts produced
    shirts_produced = order_size

    # 12) Cartons
    cartons_needed = math.ceil(shirts_produced / CARTON_CAPACITY)
    carton_weight = SHIRT_WEIGHT * CARTON_CAPACITY

    # 13) Containers
    containers_needed = math.ceil(cartons_needed / CONTAINER_CAPACITY)
    container_weight = carton_weight * CONTAINER_CAPACITY

    return {
        "order_size": order_size,
        "cotton_bales": {
            "needed": bales_needed,
            "weight": actual_cotton_weight,
            "lots": bale_lots,
        },
        "yarn_cones": {
            "needed": cones_needed,
            "weight": YARN_CONE_WEIGHT,
            "lots": yarn_lots,
        },
        "unfinished_fabric": {
            "pieces": unfinished_fabric_pieces,
            "weight": UNFINISHED_FABRIC_WEIGHT,
            "lots": unfinished_fabric_lots,
        },
        "finished_fabric": {
            "pieces": finished_fabric_pieces,
            "weight": FINISHED_FABRIC_WEIGHT,
            "lots": finished_fabric_lots,
        },
        "cut_parts": {"quantities": cut_parts, "weight_each": cut_part_weight},
        "buttons": buttons_needed,
        "shirts_produced": shirts_produced,
        "cartons": {"needed": cartons_needed, "weight": carton_weight},
        "containers": {"needed": containers_needed, "weight": container_weight},
    }


# Calculate for each order size
order_sizes = [10000, 15000, 20000]
results = {size: calculate_supply_chain(size) for size in order_sizes}

# Print results
for size, data in results.items():
    print(f"\nOrder Size: {size}")
    print(
        f"1) Cotton Bales: {data['cotton_bales']['needed']} needed, total weight: {data['cotton_bales']['weight']:.2f} lbs"
    )
    print(f"2) Cotton Bale Lots: {data['cotton_bales']['lots']}")
    print(
        f"3) Yarn Cones: {data['yarn_cones']['needed']} needed, weight each: {data['yarn_cones']['weight']} lbs"
    )
    print(f"4) Yarn Lots: {data['yarn_cones']['lots']}")
    print(
        f"5) Unfinished Fabric: {data['unfinished_fabric']['pieces']} pieces needed, weight each: {data['unfinished_fabric']['weight']} lbs"
    )
    print(f"6) Unfinished Fabric Lots: {data['unfinished_fabric']['lots']}")
    print(
        f"7) Finished Fabric: {data['finished_fabric']['pieces']} pieces needed, weight each: {data['finished_fabric']['weight']} lbs"
    )
    print(f"8) Finished Fabric Lots: {data['finished_fabric']['lots']}")
    print(
        f"9) Cut Parts: {data['cut_parts']['quantities']}, weight each: {data['cut_parts']['weight_each']:.4f} lbs"
    )
    print(f"10) Buttons: {data['buttons']} needed")
    print(f"11) Shirts Produced: {data['shirts_produced']}")
    print(
        f"12) Cartons: {data['cartons']['needed']} needed, weight each: {data['cartons']['weight']:.2f} lbs"
    )
    print(
        f"13) Containers: {data['containers']['needed']} needed, weight each: {data['containers']['weight']:.2f} lbs"
    )
