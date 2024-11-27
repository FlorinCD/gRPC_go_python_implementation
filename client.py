import grpc
import user_pb2
import user_pb2_grpc


def register_user(stub, id, name, email):
    # create an user object
    user = user_pb2.User(id=id, name=name, email=email)
    # register the user
    response = stub.RegisterUser(user)
    print(f"Registered User ID: {response.id}")


def get_user(stub, id):
    # get user id
    user_id = user_pb2.UserID(id=id)
    # get the user based on the id
    user = stub.GetUser(user_id)
    print(f"User Retrieved: {user}")


def delete_user(stub, id):
    # get user id
    user_id = user_pb2.UserID(id=id)
    # delete the user
    stub.DeleteUser(user_id)
    print("User Deleted")


def list_users(stub):
    # list all users
    for user in stub.ListUsers(user_pb2.Empty()):
        print(f"User: {user}")


def main():
    channel = grpc.insecure_channel('localhost:50051')
    stub = user_pb2_grpc.UserServiceStub(channel)

    while True:
        print("\nUser Management System")
        print("1. Register User")
        print("2. Get User")
        print("3. List Users")
        print("4. Delete User")
        print("5. Exit")
        
        choice = input("Enter your choice (1-5): ")
        
        try:
            if choice == "1":
                name = input("Enter name: ")
                email = input("Enter email: ")
                id = input("Enter id: ")
                register_user(stub, id, name, email)
                
            elif choice == "2":
                user_id = input("Enter user ID: ")
                get_user(stub, user_id)
                
            elif choice == "3":
                list_users(stub)
                    
            elif choice == "4":
                user_id = input("Enter user ID to delete: ")
                delete_user(stub, user_id)
                
            elif choice == "5":
                print("Exiting...")
                break
                
            else:
                print("Invalid choice. Please try again.")
                
        except grpc.RpcError as e:
            print(f"Error: {e.details()}")


if __name__ == "__main__":
    main()
