%tusk.string = type { i8*, i32 }
%tuskclass.Main = type { i32, i32 }

@.Main_b = global i32 0
@.Main_c = global i32 0
@.Main_main = global void ()* null

define %tusk.string @tusk.newstring(i8* %sptr, i32 %slen) {
0:
	%1 = alloca %tusk.string
	%2 = getelementptr %tusk.string, %tusk.string* %1, i32 0, i32 0
	store i8* %sptr, i8** %2
	%3 = getelementptr %tusk.string, %tusk.string* %1, i32 0, i32 1
	store i32 %slen, i32* %3
	%4 = load %tusk.string, %tusk.string* %1
	ret %tusk.string %4
}

define void @_tusk_init() {
0:
	store void ()* @tv_2, void ()** @.Main_main
	ret void
}

define %tuskclass.Main @tuskclass.new.Main() {
0:
	%1 = alloca %tuskclass.Main
	store i32 2, i32* @.Main_b
	store i32 33, i32* @.Main_c
	%2 = load %tuskclass.Main, %tuskclass.Main* %1
	ret %tuskclass.Main %2
}

declare void @tv_1()

define void @tv_2() {
0:
	%1 = call %tuskclass.Main @tuskclass.new.Main()
	%2 = alloca %tuskclass.Main
	store %tuskclass.Main %1, %tuskclass.Main* %2
	%3 = getelementptr %tuskclass.Main, %tuskclass.Main* %2, i32 0, i32 0
	store i32 3, i32* %3
	%4 = getelementptr %tuskclass.Main, %tuskclass.Main* %2, i32 0, i32 1
	store i32 33, i32* %4
	ret void
}

define void @main() {
0:
	call void @_tusk_init()
	ret void
}
