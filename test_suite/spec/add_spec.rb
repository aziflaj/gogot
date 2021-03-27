RSpec.describe 'Adding files' do
  let(:indexes) { IO.readlines('.gogot/index') }

  around do |example|
    `gogot init kewl-projekt`
    Dir.chdir('kewl-projekt') do
      example.run
    end
    `rm -rf kewl-projekt`
  end

  context 'without a path' do
    let(:command) { `gogot add` }

    it 'prints usage information' do
      expect(command).to include('Usage: gogot add <path1> [<path2>] ...')
    end
  end

  describe 'adding a file' do
    let(:command) { `gogot add ./#{filename}` }
    let(:filename) { 'hello.txt' }

    before do
      File.open(filename, 'w+') { |f| f.write("Howdy y'all") }

      command
    end

    it 'creates an index' do
      expect(indexes.count).to eq(1)
      _, filepath = indexes[0].split

      expect(filepath).to eq("./#{filename}")
    end

    it 'generates blob' do
      expect(indexes.count).to eq(1)
      hash, = indexes[0].split

      expect(File).to exist(".gogot/objects/#{hash[0..1]}/#{hash[2..]}")
    end
  end

  describe 'adding a directory' do
    before do
      File.open('file1.txt', 'w+') { |f| f.write('Test File 1') }

      Dir.mkdir('testdir')
      File.open('testdir/file11.txt', 'w+') { |f| f.write('Test File 11 in testdir') }
      File.open('testdir/file12.txt', 'w+') { |f| f.write('Test File 12 in testdir') }
      File.open('testdir/file13.txt', 'w+') { |f| f.write('Test File 13 in testdir') }

      Dir.mkdir('testdir/tester-dir')
      File.open('testdir/tester-dir/file11.txt', 'w+') { |f| f.write('Test File 11 in tester-dir') }
      File.open('testdir/tester-dir/file12.txt', 'w+') { |f| f.write('Test File 12 in tester-dir') }
      File.open('testdir/tester-dir/file13.txt', 'w+') { |f| f.write('Test File 13 in tester-dir') }

      Dir.mkdir('testdir2')
      File.open('testdir2/file21.txt', 'w+') { |f| f.write('Test File 21 in testdir2') }
      File.open('testdir2/file22.txt', 'w+') { |f| f.write('Test File 22 in testdir2') }
      File.open('testdir2/file23.txt', 'w+') { |f| f.write('Test File 23 in testdir2') }

      Dir.mkdir('testdir3') # Empty dir shall be ignored
    end

    context 'with ignore file' do
      before do
        File.open('.gogotignore', 'w+') do |f|
          f.puts './testdir2'
          f.puts './testdir/tester-dir/file11.txt'
        end
      end

      it 'excludes directories from index' do
        `gogot add .`

        expect(indexes.join("\n")).not_to include('testdir2')
      end

      it 'excludes files from index' do
        `gogot add .`

        expect(indexes.join("\n")).not_to include('testdir/tester-dir/file11')
      end
    end

    it 'ignores empty dirs' do
      `gogot add .`

      expect(indexes.join("\n")).not_to include('testdir3')
    end

    it 'adds dirs recursively to index' do
      `gogot add .`

      expect(indexes.count).to eq(1 + 3 + 3 + 3)
    end

    it 'add multiple dirs recursively' do
      `gogot add ./testdir/tester-dir ./testdir2`

      expect(indexes.count).to eq(3 + 3)
    end
  end
end
