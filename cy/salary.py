def get_salary(hours: int) -> int:
    res = hours * 84
    if hours > 120:
        res += (hours - 120) * 84 * 0.15
    if hours < 60:
        res -= 700
    return res


if __name__ == "__main__":
    try:
        hours = int(input("Pls input the hours of work:"))
        print(f"Salary is {get_salary(hours)}")
    except Exception as e:
        print("The type of input is not integer,and the error is ", e)

